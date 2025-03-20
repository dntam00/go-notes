package main

import (
	"log"
)

// https://medium.com/@philbrainy/go-slices-demystified-a-deep-dive-into-memory-layout-and-behavior-59cffd1a49ca

func main() {
	s := make([]int, 5, 5)
	s[0] = 1
	log.Println(s)
	log.Printf("Before: len=%d, cap=%d, ptr=%p, sAddr=%p\n", len(s), cap(s), s, &s)
	newArray := append(s, 1)

	s[3] = 5
	log.Println(s)

	log.Println(newArray)

	log.Printf("After: len=%d, cap=%d, ptr=%p, sAddr=%p\n", len(newArray), cap(newArray), newArray, &newArray)
}
