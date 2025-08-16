package main

import "math"

func main() {
	size := 4
	maxSize := math.MaxUint32

	println((maxSize - 1) & (size - 1))
	println(maxSize & (size - 1))
	println((maxSize + 1) & (size - 1))
	println((maxSize + 2) & (size - 1))
}
