package main

import (
	"fmt"
	"unsafe"
)

func main() {
	interfaceSize()
}

func interfaceSize() {
	var i interface{}
	i = int64(5)
	iPtr := &i
	iUintPtr := unsafe.Pointer(iPtr)
	// Result: 8 bytes
	// The size of pointer doesn't depend on the type being pointed to, only on the system architecture.
	fmt.Println(unsafe.Sizeof(iUintPtr))
	// Result: 16 bytes
	fmt.Println(unsafe.Sizeof(i))
}

func intSize() {
	x := int64(5)
	xPtr := &x
	xUintPtr := unsafe.Pointer(xPtr)
	fmt.Println(unsafe.Sizeof(xUintPtr))
}
