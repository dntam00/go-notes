package main

import (
	"context"
	"github.com/redis/rueidis"
	"log"
	"time"
)

func main() {
	//dedicateConn()
	shareConn()
}

func dedicateConn() {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"127.0.0.1:6379"},
		Password:    "turnserver",
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	for {
		log.Println("start pub/sub")
		conn, _ := client.Dedicate()
		wait := conn.SetPubSubHooks(rueidis.PubSubHooks{
			OnMessage: func(m rueidis.PubSubMessage) {

			},
		})
		log.Println("pub/sub is ready")
		err := <-wait
		log.Println("pub/sub error", err)
		conn.Close()
		time.Sleep(1 * time.Millisecond)
	}
}

func shareConn() {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"127.0.0.1:6379"},
		Password:    "turnserver",
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	i := 0

	go func() {
		for {
			if i == 10 {
				cancelFunc()
				break
			}
			time.Sleep(10 * time.Second)
		}
	}()
	for {
		err = client.Receive(ctx, client.B().Subscribe().Channel("test:channel").Build(), func(msg rueidis.PubSubMessage) {
			log.Println("Receive message:", msg.Message)
			i++
		})

		log.Println("release pub/sub", err)

		time.Sleep(1 * time.Millisecond)
	}

}
