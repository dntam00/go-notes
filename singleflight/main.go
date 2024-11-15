package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/singleflight"
	"math/rand"
	"os/signal"
	"syscall"
	"time"
)

var (
	group = &singleflight.Group{}
	n     = 10
)

func do() int {
	r := rand.Intn(1000) + 1
	time.Sleep(time.Duration(r) * time.Millisecond)
	return r
}

func main() {
	key := "flight"
	for i := 0; i < n; i++ {
		go func() {
			v, _, _ := group.Do(key, func() (interface{}, error) {
				return do(), nil
			})
			fmt.Println("result: ", v)
		}()
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
}
