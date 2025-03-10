package main

import (
	"bufio"
	"fmt"
	"log"
	"sync/atomic"
	"testing"
)

func BenchmarkWriting(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := &dumpWriter{
			count: atomic.Int64{},
		}
		input := make(chan []byte, 1)
		buf := bufio.NewWriterSize(w, 120)
		log.Println("start program")
		closeCh := make(chan struct{})
		go func() {
			err := writing(buf, input)
			if err != nil {
				log.Fatalf("Error writing: %v", err)
			}
			closeCh <- struct{}{}
		}()

		go func() {
			for i := 0; i < 100000000; i++ {
				input <- []byte("hello world ")
			}
			log.Println("finish send")
			close(input)
		}()

		fmt.Println("finish test")

		select {
		case <-closeCh:
		}
	}
}
