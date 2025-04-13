package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	var num int64 = -1

	fmt.Printf("- Binary representation of %d in varint: ", num)

	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(num))
	for _, b := range buf {
		fmt.Printf("%08b ", b)
	}
	fmt.Println()

	fmt.Printf("- Binary representation of %d in zigzag varint: ", num)
	bufZigzag := make([]byte, 1)
	binary.PutVarint(bufZigzag, num)
	for _, b := range bufZigzag {
		fmt.Printf("%08b ", b)
	}
}
