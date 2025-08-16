package common

import (
	"fmt"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislock"
	"net"
	"strings"
	"time"
)

func InitRedis() rueidis.Client {
	addresses := strings.Split("127.0.0.1:7500", ",")

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
	_ = client.B().Ping().Build()
	fmt.Println("Connected to redis:", addresses)
	return client
}

func InitWithLock() (rueidis.Client, rueidislock.Locker) {
	address := "127.0.0.1:6379"
	addresses := strings.Split(address, ",")

	clientOption := rueidis.ClientOption{
		Dialer: net.Dialer{
			Timeout:   time.Millisecond * time.Duration(60000),
			KeepAlive: time.Millisecond * time.Duration(10000),
		},
		InitAddress: addresses,
		Username:    "",
		Password:    "turnserver",
		SelectDB:    0,
	}
	client, err := rueidis.NewClient(clientOption)
	if err != nil {
		panic(err)
	}
	_ = client.B().Ping().Build()
	fmt.Println("connected to redis server:", address)

	lockerOption := rueidislock.LockerOption{
		KeyPrefix:      "kaixin",
		ClientOption:   clientOption,
		KeyValidity:    time.Second * 5,
		ExtendInterval: time.Second * 1,
		TryNextAfter:   time.Millisecond * 200,
	}
	locker, _ := rueidislock.NewLocker(lockerOption)

	return client, locker
}
