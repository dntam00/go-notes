package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	timeout, _ := context.WithTimeout(context.Background(), 1000000*time.Second)
	timeout = context.WithValue(timeout, "key", "value")
	time.Sleep(5 * time.Second)
	if timeout.Err() != nil {
		fmt.Println("error")
	} else {
		fmt.Println("not error")
	}
}
