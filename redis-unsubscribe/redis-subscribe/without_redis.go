package main

//
//import (
//	"context"
//	"fmt"
//	"github.com/redis/rueidis"
//	"net"
//	"os/signal"
//	"strconv"
//	"strings"
//	"sync/atomic"
//	"syscall"
//	"time"
//)
//
//func main() {
//	redis := InitRedis()
//
//	trafficCh := "turn/realm/*/user/*/allocation/*/traffic/peer"
//	statusCh := "turn/realm/*/user/*/allocation/*/status"
//
//	numsDeleted := int32(0)
//
//	go func() {
//		err := redis.Receive(context.Background(), redis.B().Psubscribe().Pattern(statusCh).Build(), func(msg rueidis.PubSubMessage) {
//			if msg.Message == "deleted" {
//				atomic.AddInt32(&numsDeleted, 1)
//			}
//		})
//		if err != nil {
//			fmt.Println("Error subscribing to channel", statusCh)
//		}
//	}()
//
//	nums := int32(0)
//
//	go func() {
//		err := redis.Receive(context.Background(), redis.B().Psubscribe().Pattern(trafficCh).Build(), func(msg rueidis.PubSubMessage) {
//			atomic.AddInt32(&nums, 1)
//		})
//		if err != nil {
//			fmt.Println("Error subscribing to channel", trafficCh)
//		}
//	}()
//
//	go func() {
//		ticker := time.NewTicker(60 * time.Second)
//		for {
//			select {
//			case <-ticker.C:
//				fmt.Println(fmt.Sprintf("Number of messages received: status: %d, traffic %d", atomic.LoadInt32(&numsDeleted), atomic.LoadInt32(&nums)))
//			}
//		}
//	}()
//
//	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
//	defer stop()
//	<-ctx.Done()
//}
//
//type PeerTrafficMessage struct {
//	Rcvp  int `json:"rcvp"`
//	Rcvb  int `json:"rcvb"`
//	Sentp int `json:"sentp"`
//	Sentb int `json:"sentb"`
//}
//
//func ConvertToPeerTrafficMessage(message string) (PeerTrafficMessage, error) {
//	msg := strings.Split(message, ", ")
//	rcvp, _ := strconv.Atoi(strings.Split(msg[0], "=")[1])
//	rcvb, _ := strconv.Atoi(strings.Split(msg[1], "=")[1])
//	sentp, _ := strconv.Atoi(strings.Split(msg[2], "=")[1])
//	sentb, _ := strconv.Atoi(strings.Split(msg[3], "=")[1])
//
//	return PeerTrafficMessage{
//		Rcvp:  rcvp,
//		Rcvb:  rcvb,
//		Sentp: sentp,
//		Sentb: sentb,
//	}, nil
//}
//
//func GetAllocationInfo(topic string) (string, string) {
//	parts := strings.Split(topic, "/")
//	return parts[2], parts[4] + "-" + parts[6]
//}
//
//func InitRedis() rueidis.Client {
//	//addresses := strings.Split("127.0.0.1:6379", ",")
//
//	clientOption := rueidis.ClientOption{
//		Dialer: net.Dialer{
//			Timeout:   time.Millisecond * time.Duration(60000),
//			KeepAlive: time.Millisecond * time.Duration(10000),
//		},
//		InitAddress:  addresses,
//		Username:     "",
//		Password:     "",
//		SelectDB:     0,
//		DisableRetry: true,
//	}
//	client, err := rueidis.NewClient(clientOption)
//	if err != nil {
//		panic(err)
//	}
//	_ = client.B().Ping().Build()
//	fmt.Println("Connected to redis, start program")
//	return client
//}
