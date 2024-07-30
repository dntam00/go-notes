package main

import (
	"context"
	"log"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	withValue := context.WithValue(ctx, "key", "value")
	cancelFunc()
	value := withValue.Value("key")
	log.Println(value)
	cancelFunc()
	cancelFunc()
	cancelFunc()
	log.Println(withValue.Value("key"))
}
