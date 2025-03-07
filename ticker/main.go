package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ticker := time.NewTicker(2 * time.Second)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	for range ticker.C {
		select {
		case <-ctx.Done():
			ticker.Stop()
			fmt.Println("Ticker stopped")
			return
		default:
		}
		fmt.Println("Tick at", time.Now())
		time.Sleep(3 * time.Second)
	}
}
