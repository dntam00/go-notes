package main

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

// convert int to byte chunks.
func main() {
	n := 257
	formatInt := strconv.FormatInt(int64(n), 2)
	first := formatInt[:len(formatInt)-8]
	second := formatInt[len(formatInt)-8:]

	i, _ := strconv.ParseInt(first, 2, 64)
	y, _ := strconv.ParseInt(second, 2, 64)

	fmt.Println(i)
	fmt.Println(y)

	fmt.Println(binary.LittleEndian.Uint64([]byte(first)))
	fmt.Println(formatInt)
}
