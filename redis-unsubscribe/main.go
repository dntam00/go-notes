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
	//go connect("abc")
	go connectInDedicated("abc")
	time.Sleep(1 * time.Second)
	//go connect("xyz")
	go connectInDedicated("xyz")
	time.Sleep(10 * time.Second)
	go func() {
		redis.Do(context.Background(), redis.B().Unsubscribe().Channel("abc").Build())
	}()
	go func() {
		redis.Do(context.Background(), redis.B().Unsubscribe().Channel("xyz").Build())
	}()
	time.Sleep(50 * time.Second)
	fmt.Println("End")
}

func connect(channel string) {
	err := redis.Receive(context.Background(), redis.B().Subscribe().Channel(channel).Build(), func(msg rueidis.PubSubMessage) {
		fmt.Println("Receive message:", msg.Message)
	})
	if err != nil {
		fmt.Println("Error when subscribing channel:", err)
	}
	fmt.Println("Release connection")
}

func connectInDedicated(channel string) {
	err := redis.Dedicated(func(client rueidis.DedicatedClient) error {
		return client.Receive(context.Background(), client.B().Subscribe().Channel(channel).Build(), func(msg rueidis.PubSubMessage) {
			fmt.Println("Receive message:", msg.Message)
		})
	})
	fmt.Println(err)
}
