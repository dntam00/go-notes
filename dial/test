package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

func main() {
	ch := make(chan struct{}, 5)
	ch <- struct{}{}
	m := <-ch
	fmt.Println(m)

	dialer := net.Dialer{
		Timeout:   time.Millisecond * time.Duration(10000),
		KeepAlive: time.Millisecond * time.Duration(10000),
	}

	_, err := dialer.DialContext(context.Background(), "tcp", "0.0.0.0:6379")
	if err != nil {
		panic(err)
	}
	fmt.Println("success")
}
