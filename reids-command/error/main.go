package main

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislock"
	"play-around/common"
	"time"
)

var client rueidis.Client
var locker rueidislock.Locker

func main() {
	client, locker = common.InitWithLock()
	//resp := setCmd()
	//fmt.Println(resp.Error())
	lockCmd()
}

func setCmd() rueidis.RedisResult {
	return client.Do(context.Background(), client.B().Set().Key("key").Value("value").Pxat(time.Now().Add(time.Duration(1000000)*time.Millisecond)).Build())
}

func lockCmd() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFunc()
	_, _, err := locker.WithContext(ctx, "key")
	//_, _, err := locker.TryWithContext(ctx, "key")
	fmt.Println(err)
}
