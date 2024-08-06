package main

import (
	"context"
	"log"
	"time"
)

func main() {
	timeout, _ := context.WithTimeout(context.Background(), 5*time.Second)
	ctx, cancelFunc := context.WithCancel(timeout)
	time.Sleep(6 * time.Second)
	if err := timeout.Err(); err != nil {
		log.Println(err)
	}
	withValue := context.WithValue(ctx, "key", "value")
	cancelFunc()
	value := withValue.Value("key")
	log.Println(value)
	cancelFunc()
	cancelFunc()
	cancelFunc()
	log.Println(withValue.Value("key"))
}
