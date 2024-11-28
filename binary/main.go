package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

// convert int to byte chunks.
func main() {
	convertToBytes()
}

func convert() {
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

func sizeBytes() {
	bytes := [3000]byte{}
	fmt.Println(len(bytes))
}

const (
	MoreFragments int = 1 << iota // more fragments flag
	DontFragment                  // don't fragment flag
)

func ipFlags() {
	//flags := ipv4.HeaderFlags(0x40)
	//
	//formatInt := strconv.FormatInt(int64(flags), 2)
	fmt.Println(DontFragment)
}

func convertToBytes() {
	// The number to convert
	number := uint16(255) // Example: 500

	// Create a byte buffer
	buf := new(bytes.Buffer)

	// Write the number to the buffer in big-endian format
	err := binary.Write(buf, binary.BigEndian, number)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Get the 2-byte representation
	twoBytes := buf.Bytes()

	// Print the result
	fmt.Printf("Number: %d, Two Bytes: %v\n", number, twoBytes)
}
