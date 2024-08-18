package main

import (
	"context"
	"fmt"
	"play-around/common"
	"time"
)

func main() {
	_, locker := common.InitWithLock()
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	_, _, err := locker.WithContext(ctx, "test")
	if err != nil {
		fmt.Println(err)
	}
}
