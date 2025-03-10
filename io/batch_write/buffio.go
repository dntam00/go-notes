package main

import (
	"bufio"
	"log"
	"sync/atomic"
)

type dumpWriter struct {
	count atomic.Int64
}

func (w *dumpWriter) Write(p []byte) (n int, err error) {
	w.count.Add(1)
	return len(p), nil
}

func main() {
	//w := os.Stdout
	//w := &dumpWriter{
	//	count: atomic.Int64{},
	//}
	w := log.Writer()
	input := make(chan []byte, 10000)
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
		for i := 0; i < 10000; i++ {
			input <- []byte("hello world ")
		}
		log.Println("finish send")
		close(input)
	}()

	select {
	case <-closeCh:
	}
	log.Println("\nend program", counter.Load())
}

var counter = atomic.Int64{}

func writing(buf *bufio.Writer, input <-chan []byte) (err error) {
	var data []byte
	var more = true
	for more && err == nil {
		select {
		case data, more = <-input:
			//log.Println("input case")
			counter.Add(1)
			_, err = buf.Write(data)
			continue
		default:
			//log.Println("default case")
		}
		if err = buf.Flush(); err == nil {
			//log.Println("flush")
			data, more = <-input
			_, err = buf.Write(data)
		}
	}
	_ = buf.Flush()
	return err
}
