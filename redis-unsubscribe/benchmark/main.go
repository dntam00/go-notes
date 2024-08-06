package main

import (
	"context"
	"flag"
	"github.com/redis/rueidis"
	"log"
	"net/http"
	_ "net/http/pprof"
	"play-around/common"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

var client rueidis.Client

var numNodes = flag.Int("nodes", 200, "num nodes")

var (
	read        int64
	write       int64
	concurrency = runtime.NumCPU()
	clients     = 500
	data        = string(make([]byte, 512))
)

func main() {
	client = common.InitRedis()
	defer client.Close()

	go func() {
		go func() {
			log.Println(http.ListenAndServe("localhost:5001", nil))
		}()
	}()

	go func() {
		var valRead int64
		var valWrite int64
		for {
			time.Sleep(time.Second)
			newValRead := atomic.LoadInt64(&read)
			newValWrite := atomic.LoadInt64(&write)
			log.Println(newValRead-valRead, "msg/sec read", newValWrite-valWrite, "msg/sec write")
			valWrite = newValWrite
			valRead = newValRead
		}
	}()
	run()
	//runInDedicated()
}

func run() {

	messages := make(chan rueidis.PubSubMessage, 10000)

	for i := 0; i < clients; i++ {
		go func(index int) {
			var c string
			c = "request-" + strconv.Itoa(index)
			err := client.Receive(context.Background(), client.B().Subscribe().Channel(c, "response-"+strconv.Itoa(index)).Build(), func(msg rueidis.PubSubMessage) {
				atomic.AddInt64(&read, 1)
				messages <- msg
			})
			if err != nil {
				panic(err)
			}
		}(i)
	}

	for z := 0; z < concurrency; z++ {
		go func() {
			for {
				for i := 0; i < clients; i++ {
					//var c string
					//c = "request-" + strconv.Itoa(i)
					//redis.Do(context.Background(), redis.B().Publish().Channel(c).Message(data).Build())
					//atomic.AddInt64(&write, 1)

					client.Do(context.Background(), client.B().Publish().Channel("response-"+strconv.Itoa(i)).Message(data).Build())
					atomic.AddInt64(&write, 1)
					<-messages
				}
				//for i := 0; i < clients; i++ {
				//	<-messages
				//}
			}
		}()
	}
	select {}
}

func runInDedicated() {
	defer client.Close()

	messages := make(chan rueidis.PubSubMessage, 1024)
	var subscribeConns []rueidis.DedicatedClient

	for i := 0; i < *numNodes; i++ {
		conn, _ := client.Dedicate()

		conn.SetPubSubHooks(rueidis.PubSubHooks{
			OnMessage: func(m rueidis.PubSubMessage) {
				atomic.AddInt64(&read, 1)
				messages <- m
			},
		})

		resp := conn.Do(context.Background(), client.B().Subscribe().Channel("request", "response"+strconv.Itoa(i)).Build())
		if resp.Error() != nil {
			panic(resp.Error())
		}
		subscribeConns = append(subscribeConns, conn)
	}

	defer func() {
		for _, conn := range subscribeConns {
			conn.Close()
		}
	}()

	for i := 0; i < concurrency; i++ {
		go func() {
			for {
				client.Do(context.Background(), client.B().Publish().Channel("request").Message(data).Build())
				atomic.AddInt64(&write, 1)
				for i := 0; i < *numNodes; i++ {
					<-messages
					client.Do(context.Background(), client.B().Publish().Channel("response"+strconv.Itoa(i)).Message(data).Build())
					atomic.AddInt64(&write, 1)
				}
				for i := 0; i < *numNodes; i++ {
					<-messages
				}
			}
		}()
	}
	select {}
}
