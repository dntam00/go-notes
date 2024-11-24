package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	l := log.New(os.Stderr, "test ", 1)
	l.Println("log message")
	fmt.Println("fmt message")
}
