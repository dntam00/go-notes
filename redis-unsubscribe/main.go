package main

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"os/signal"
	"play-around/common"
	"syscall"
	"time"
)

var redis rueidis.Client

func main() {
	redis = common.InitRedis()
	go connect("kaixin:test")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
}

func connect(channel string) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Second)
		cancelFunc()
	}()
	err := redis.Receive(ctx, redis.B().Ssubscribe().Channel(channel).Build(), func(msg rueidis.PubSubMessage) {
		fmt.Println("receive message:", msg.Message)
	})
	if err != nil {
		fmt.Println("error when subscribing channel:", err)
	}
	fmt.Println("release connection")
}

func connectInDedicated(channel string) {
	err := redis.Dedicated(func(client rueidis.DedicatedClient) error {
		return client.Receive(context.Background(), client.B().Subscribe().Channel(channel).Build(), func(msg rueidis.PubSubMessage) {
			fmt.Println("Receive message:", msg.Message)
		})
	})
	fmt.Println(err)
}
