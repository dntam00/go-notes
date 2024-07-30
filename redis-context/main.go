package main

import (
	"context"
	"fmt"
	"play-around/common"
	"time"
)

func main() {
	_, locker := common.InitWithLock()
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	_, _, err := locker.WithContext(ctx, "lock:123")
	if err != nil {
		panic(err)
	}
	fmt.Println("success")
}
