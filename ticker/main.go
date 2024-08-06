package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			fmt.Println("Tick")
		}
	}()

	time.Sleep(time.Second * 10)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}
