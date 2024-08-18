package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	for timer := time.NewTimer(2 * time.Second); true; {
		select {
		case <-timer.C:
			fmt.Println("timer")
			timer.Reset(2 * time.Second)
		case <-ch:

		}
	}
}
