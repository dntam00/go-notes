package main

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	redis := InitRedis()

	trafficCh := "turn/realm/*/user/*/allocation/*/traffic/peer"
	statusCh := "turn/realm/*/user/*/allocation/*/status"

	statusMsgHash := "status:count:%s"
	numsMsgHash := "traffic:count:%s"

	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	go func() {
		err := redis.Receive(context.Background(), redis.B().Psubscribe().Pattern(statusCh).Build(), func(msg rueidis.PubSubMessage) {
			//go func() {
			if msg.Message == "deleted" {
				turnIp, allocationId := GetAllocationInfo(msg.Channel)
				msgHashKey := fmt.Sprintf(statusMsgHash, turnIp)
				timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancelFunc()
				now := time.Now()
				if err := redis.Do(timeout, redis.B().Hset().Key(msgHashKey).FieldValue().FieldValue(allocationId, "1").Build()).Error(); err != nil {
					fmt.Println("Error setting hash value to status", allocationId, err, now.Unix(), time.Now().Unix())
					_ = pprof.Lookup("goroutine").WriteTo(os.Stdout, 2)
				}
			}
			//}()
		})
		if err != nil {
			fmt.Println("Error subscribing to channel", statusCh)
		}
	}()

	go func() {
		err := redis.Receive(context.Background(), redis.B().Psubscribe().Pattern(trafficCh).Build(), func(msg rueidis.PubSubMessage) {
			//go func() {
			message, _ := ConvertToPeerTrafficMessage(msg.Message)
			turnIp, allocationId := GetAllocationInfo(msg.Channel)
			//sid := uuid.New().String()
			//fmt.Println("start process message", sid)
			if message.Rcvp+message.Sentp > 0 {
				timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancelFunc()
				if err := redis.Do(timeout, redis.B().Hset().Key(fmt.Sprintf(numsMsgHash, turnIp)).FieldValue().FieldValue(allocationId, strconv.Itoa(1)).Build()).Error(); err != nil {
					_ = pprof.Lookup("goroutine").WriteTo(os.Stdout, 2)
				}
			}
			//}()
		})
		if err != nil {
			fmt.Println("Error subscribing to channel", trafficCh)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
}

type PeerTrafficMessage struct {
	Rcvp  int `json:"rcvp"`
	Rcvb  int `json:"rcvb"`
	Sentp int `json:"sentp"`
	Sentb int `json:"sentb"`
}

func ConvertToPeerTrafficMessage(message string) (PeerTrafficMessage, error) {
	msg := strings.Split(message, ", ")
	rcvp, _ := strconv.Atoi(strings.Split(msg[0], "=")[1])
	rcvb, _ := strconv.Atoi(strings.Split(msg[1], "=")[1])
	sentp, _ := strconv.Atoi(strings.Split(msg[2], "=")[1])
	sentb, _ := strconv.Atoi(strings.Split(msg[3], "=")[1])

	return PeerTrafficMessage{
		Rcvp:  rcvp,
		Rcvb:  rcvb,
		Sentp: sentp,
		Sentb: sentb,
	}, nil
}

func GetAllocationInfo(topic string) (string, string) {
	parts := strings.Split(topic, "/")
	return parts[2], parts[4] + "-" + parts[6]
}

func InitRedis() rueidis.Client {
	addresses := strings.Split("127.0.0.1:6379", ",")

	clientOption := rueidis.ClientOption{
		Dialer: net.Dialer{
			Timeout:   time.Millisecond * time.Duration(60000),
			KeepAlive: time.Millisecond * time.Duration(10000),
		},
		InitAddress:  addresses,
		Username:     "",
		Password:     "turnserver",
		SelectDB:     0,
		DisableRetry: true,
	}
	client, err := rueidis.NewClient(clientOption)
	if err != nil {
		panic(err)
	}
	_ = client.B().Ping().Build()
	fmt.Println("Connected to Redis")
	return client
}
