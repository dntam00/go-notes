package main

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"play-around/common"
	"strings"
	"time"
)

func main() {
	wrongMessage := 0
	allocations := make(map[string]time.Time)
	channel := "turn/realm/*/user/*/allocation/*/traffic/peer"
	redis := common.InitRedis()
	err := redis.Receive(context.Background(), redis.B().Psubscribe().Pattern(channel).Build(), func(msg rueidis.PubSubMessage) {
		allocationId := getAllocationInfo(msg.Channel)
		allocation, ok := allocations[allocationId]
		if !ok {
			allocations[allocationId] = time.Now()
		} else if time.Now().Before(allocation.Add(14 * time.Second)) {
			wrongMessage++
		}
		if wrongMessage > 0 {
			panic("Too many wrong messages")
		}
		fmt.Println("receive message", msg.Message)
	})
	if err != nil {
		fmt.Println("Error when subscribing channel:", err)
	}
	fmt.Println("Release connection")
}

func getAllocationInfo(topic string) string {
	parts := strings.Split(topic, "/")
	return parts[6]
}
