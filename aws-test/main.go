package main

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"net"
	"strings"
	"time"
)

func InitRedis() rueidis.Client {
	addresses := strings.Split("", ",")

	clientOption := rueidis.ClientOption{
		Dialer: net.Dialer{
			Timeout:   time.Millisecond * time.Duration(60000),
			KeepAlive: time.Millisecond * time.Duration(10000),
		},
		InitAddress: addresses,
		Username:    "",
		Password:    "",
		SelectDB:    0,
	}
	client, err := rueidis.NewClient(clientOption)
	if err != nil {
		panic(err)
	}
	pong := client.B().Ping().Build()
	fmt.Println("Connected to Redis:", pong)
	time.Sleep(2 * time.Second)
	fmt.Println("Number of nodes: ", client.Nodes())
	return client
}

func main() {
	redis := InitRedis()

	ctx, cancelFunc := context.WithCancel(context.Background())

	go func() {
		cmd1 := redis.B().Subscribe().Channel("channel-1").Build()
		fmt.Println("slot 1: ", cmd1.Slot())
		err := redis.Receive(ctx, cmd1, func(msg rueidis.PubSubMessage) {
			fmt.Println("1st Receive message:", msg.Message)
		})
		if err != nil {
			fmt.Println("Error when subscribing channel:", err)
		}
		fmt.Println("Release connection channel-1 1st")
	}()

	go func() {
		cmd2 := redis.B().Ssubscribe().Channel("channel-1").Build()
		fmt.Println("slot 1: ", cmd2.Slot())
		err := redis.Receive(context.Background(), cmd2, func(msg rueidis.PubSubMessage) {
			fmt.Println("2nd Receive message:", msg.Message)
		})
		if err != nil {
			fmt.Println("Error when subscribing channel:", err)
		}
		fmt.Println("Release connection channel-1 2nd")
	}()

	time.Sleep(3 * time.Second)

	go func() {
		cmd3 := redis.B().Publish().Channel("channel-1").Message("Hello").Build()
		fmt.Println("send slot: ", cmd3.Slot())
		if err := redis.Do(context.Background(), cmd3).Error(); err != nil {
			fmt.Println(err)
		}
	}()

	cancelFunc()

	cmd := redis.B().Sunsubscribe().Channel("channel-1").Build()
	fmt.Println("slot: ", cmd.Slot())
	if err := redis.Do(context.Background(), cmd).Error(); err != nil {
		fmt.Println(err)
	}
	time.Sleep(5 * time.Second)
}
