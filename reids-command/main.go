package main

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"play-around/common"
	"time"
)

var client rueidis.Client

func main() {
	redis := common.InitRedis()
	client = redis
	if err := SetDataOnlyIfExistWithTTL(context.Background(), "key", "value", 1000000); err != nil {
		fmt.Print(err)
	}
}

func SetDataOnlyIfExistWithTTL(ctx context.Context, key string, value string, ttlInMilliseconds int) error {
	return client.Do(ctx, client.B().Set().Key(key).Value(value).Xx().Pxat(time.Now().Add(time.Duration(ttlInMilliseconds)*time.Millisecond)).Build()).Error()
}
