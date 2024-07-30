package main

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"play-around/common"
	"time"
)

var redis rueidis.Client

func main() {
	redis = common.InitRedis()
	go connect()
	go connect()
	time.Sleep(10 * time.Second)
	go func() {
		redis.Do(context.Background(), redis.B().Unsubscribe().Channel("test:channel").Build())
	}()
	time.Sleep(50 * time.Second)
	fmt.Println("End")
}

func connect() {
	err := redis.Receive(context.Background(), redis.B().Subscribe().Channel("test:channel").Build(), func(msg rueidis.PubSubMessage) {
		fmt.Println("Receive message:", msg.Message)
	})
	if err != nil {
		fmt.Println("Error when subscribing channel:", err)
	}
	fmt.Println("Release connection")
}
