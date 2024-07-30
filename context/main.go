package main

import (
	"context"
	"time"
)

func main() {
	background := context.Background()
	go func() {
		time.Sleep(time.Second * 2)
		background.Done()
	}()
	timeout, cancelFunc := context.WithTimeout(background, time.Second*10)
	timeout2, _ := context.WithTimeout(timeout, time.Second*10)
	go func() {
		time.Sleep(time.Second * 3)
		cancelFunc()
	}()

	time.Sleep(4 * time.Second)

	select {
	case <-timeout2.Done():
		println("timeout: ", timeout2.Err())
		return
	default:
		println("not timeout")
	}
}
