package main

import (
	"log"
	"unsafe"
)

func main() {
	myString := "neato burrito"
	byteSlice := unsafe.Slice(unsafe.StringData(myString), len(myString))
	log.Println(byteSlice) // [110 101 97 116 111 32 98 117 114 114 105 116 111]

	//bytes := []byte(myString)

	if byteSlice[2]-'a' == 0 {
		log.Println("check")
	}

	// ByteSliceToString
	//myBytes := []byte{
	//	115, 111, 32, 109, 97, 110, 121, 32, 110,
	//	101, 97, 116, 32, 98, 121, 116, 101, 115,
	//}
}
