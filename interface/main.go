package main

import "fmt"

type conn interface {
	Do()
}

func main() {
	fmt.Println(get())
}

func get() (p conn) {
	return p
}
