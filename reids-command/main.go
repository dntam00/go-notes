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
	//Decrby()
	IncrbyTTL(context.Background(), "test", 1, 10000)
	//if err := SetDataOnlyIfExistWithTTL(context.Background(), "key", "value", 1000000); err != nil {
	//	fmt.Print(err)
	//}
	//HSET()
	//PExpireAt()
	//Get()
}

func SetDataOnlyIfExistWithTTL(ctx context.Context, key string, value string, ttlInMilliseconds int) error {
	return client.Do(ctx, client.B().Set().Key(key).Value(value).Xx().Pxat(time.Now().Add(time.Duration(ttlInMilliseconds)*time.Millisecond)).Build()).Error()
}

func HSET() {
	client.Do(context.Background(), client.B().Hset().Key("test").FieldValue().FieldValue("abc", "value123").Build())
}

func PExpireAt() {
	for {
		resp := client.Do(context.Background(), client.B().Pexpireat().Key("test").MillisecondsTimestamp(time.Now().Add(time.Duration(40000)*time.Millisecond).UnixMilli()).Build())
		v, err := resp.AsInt64()
		if err != nil {
			fmt.Println("error: ", err)
		}
		if v == 0 {
			fmt.Println("not found")
		}
		time.Sleep(5 * time.Second)
	}
}

func Get() {
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	time.Sleep(3 * time.Second)
	result := client.Do(timeout, client.B().Get().Key("test").Build())
	v, err := result.AsInt64()
	if err != nil {
		fmt.Println("error")
	}
	fmt.Println(v, err)
}

func Decrby() {
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	result := client.Do(timeout, client.B().Decrby().Key("test").Decrement(1).Build())
	v, err := result.AsInt64()
	fmt.Println(v, err)
}

func IncrbyTTL(ctx context.Context, key string, increment int64, ttlInMillis int) {
	incrByCmd := client.B().Incrby().Key(key).Increment(increment).Build()
	expireCmd := client.B().Pexpireat().Key(key).MillisecondsTimestamp(time.Now().Add(time.Duration(ttlInMillis) * time.Millisecond).UnixMilli()).Build()
	resp := client.DoMulti(ctx, incrByCmd, expireCmd)
	v, err := resp[0].AsInt64()
	if err != nil {
		return
	}
	if resp[1].Error() != nil {
		return
	}
	fmt.Println("resp 0: ", v)
	fmt.Println("resp 1: ", resp[1].String())
}
