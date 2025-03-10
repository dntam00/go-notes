package main

import (
	"log"
)

func main() {
	ch := make(chan int, 5)

	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)

	for {
		select {
		case v, ok := <-ch:
			if !ok {
				log.Println("channel closed")
				return
			}
			log.Println("v:", v)
		}
	}
}
