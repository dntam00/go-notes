package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	l := &sync.RWMutex{}

	read := sync.NewCond(l.RLocker())
	write := sync.NewCond(l)

	x := 0

	go func() {
		for {
			l.Lock()
			if x == 1 {
				write.Wait()
			}
			fmt.Println("write")
			x = 1
			l.Unlock()
			read.Signal()
		}
	}()

	//time.Sleep(1 * time.Second)

	go func() {
		for {
			l.RLock()
			if x == 0 {
				read.Wait()
			}
			fmt.Println("read")
			x = 0
			l.RUnlock()
			write.Signal()
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
}
