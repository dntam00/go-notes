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
	addresses := strings.Split("127.0.0.1:6379", ",")

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
	pong := client.B().Ping().Build()
	fmt.Println("Connected to Redis:", pong)
	return client
}

func InitWithLock() (rueidis.Client, rueidislock.Locker) {
	addresses := strings.Split("127.0.0.1:6379", ",")

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
	pong := client.B().Ping().Build()
	fmt.Println("Connected to Redis:", pong)

	lockerOption := rueidislock.LockerOption{
		ClientOption: clientOption,
		KeyValidity:  time.Second * 100,
	}
	locker, _ := rueidislock.NewLocker(lockerOption)

	return client, locker
}
