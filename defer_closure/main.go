package main

import "fmt"

func main() {
	var i int
	defer func() {
		fmt.Print(i)
	}()
	for ; i < 5; i++ {

	}
}
