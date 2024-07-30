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

	// recover avoid crash in go routine is propagated to main
	go runWithRecover(func() {
		fmt.Println("start function normally")
		panic("panic in run with recover")
	})

	time.Sleep(5 * time.Second)

	fmt.Println("returned normally from main")
}
