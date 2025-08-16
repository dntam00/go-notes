package main

import (
	"context"
	"fmt"
	"github.com/redis/rueidis/rueidislock"
	"play-around/common"
	"play-around/utils"
	"time"
)

var locker rueidislock.Locker

func main() {
	_, locker = common.InitWithLock()
	for i := 0; i < 4000; i++ {
		go func(index int) {
			ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancelFunc()
			_, cancelLockFunc := acquire(ctx, i)
			time.Sleep(time.Second * 1)
			ctx, cancelFunc2 := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancelFunc2()
			_, cancelLockFunc2 := acquire(ctx, i)
			if cancelLockFunc != nil {
				cancelLockFunc()
			}
			if cancelLockFunc2 != nil {
				cancelLockFunc2()
			}
		}(i)
	}
	utils.Wait()
}

func acquire(ctx context.Context, index int) (context.Context, context.CancelFunc) {
	lockCtx, aqrCancel, err := locker.WithContext(ctx, fmt.Sprintf("%d", index))
	if err != nil {
		fmt.Printf("Error acquiring lock for %d: %v\n", index, err)
		return lockCtx, aqrCancel
	}
	fmt.Printf("Lock acquired for %d\n", index)
	return lockCtx, aqrCancel
}
