package main

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"play-around/common"
)

func main() {
	redis := common.InitRedis()
	resp, err := redis.Do(context.Background(), redis.B().Get().Key("yy").Build()).AsInt64()
	if rueidis.IsRedisNil(err) {
		fmt.Println("Key not found")
	}
	fmt.Print(resp)
}
