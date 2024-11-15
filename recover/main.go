package main

import (
	"fmt"
	"time"
)

func runWithRecover(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from ", r)
		}
	}()
	fn()
}

func main() {

	//now := time.Now()
	//time.Sleep(500 * time.Millisecond)
	//seconds := int64(time.Since(now).Seconds())
	//fmt.Println(seconds)

	// recover avoid crash in go routine is propagated to main
	runWithRecover(func() {
		fmt.Println("start function normally")
		time.Sleep(2 * time.Second)
		panic("panic in run with recover")
	})

	fmt.Println("returned normally from main")

	time.Sleep(5 * time.Second)

	fmt.Println("after normally from main")
}
