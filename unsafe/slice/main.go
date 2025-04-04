package main

import (
	"fmt"
	"log"
	"reflect"
	"sync/atomic"
	"unsafe"
)

// https://medium.com/@philbrainy/go-slices-demystified-a-deep-dive-into-memory-layout-and-behavior-59cffd1a49ca

func main() {
	ModifySlice()
}

func modify() {
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

func SliceAddr() {
	arr := make([]int, 2)
	log.Println("addr:", (unsafe.Pointer)(&arr))
	log.Println("addr [0]:", (unsafe.Pointer)(&arr[0]))
	log.Println("addr [1]:", (unsafe.Pointer)(&arr[1]))
	log.Println("addr converted:", (*unsafe.Pointer)((unsafe.Pointer)(&arr[0])))
}

// https://blog.devtrovert.com/p/what-is-unsafepointer-or-uintptr

func ModifySliceHeader() {
	str := "hello"
	var b []byte

	// Reinterpret the string's header as a slice header
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	// Copy the data pointer and length from the string header to the slice header
	sliceHeader.Data = strHeader.Data
	sliceHeader.Len = strHeader.Len
	sliceHeader.Cap = strHeader.Len

	fmt.Println(b)         //  [104 101 108 108 111]
	fmt.Println(string(b)) // hello
}

func ModifySlice() {
	// Example 1: Create a slice from an array pointer
	arr := [5]int64{10, 20, 30, 40, 50}
	// Get pointer to first element
	ptr := &arr[0]
	// Create a slice using that pointer with length 5
	slice := unsafe.Slice(ptr, 5)
	pointer := unsafe.Pointer(&slice[0])
	nextElem := unsafe.Pointer(uintptr(pointer) + unsafe.Sizeof(int64(0)))
	newVal := int64(0)
	atomic.StoreInt64((*int64)(nextElem), newVal)

	fmt.Println(slice) // Output: [10 20 30 40 50]

	// Example 2: Working with bytes
	//data := [8]byte{1, 2, 3, 4, 5, 6, 7, 8}
	//bytePtr := &data[0]
	//// Reinterpret as uint16 slice (assuming proper alignment)
	//uint16Slice := unsafe.Slice((*uint16)(unsafe.Pointer(bytePtr)), 4)
	//fmt.Println(uint16Slice) // Will show 4 uint16 values from the bytes
	//
	//// Example 3: Taking a slice of a single value
	//newVal := 42
	//// Create a single-element slice containing this value
	//singleSlice := unsafe.Slice(&newVal, 1)
	//fmt.Println(singleSlice) // Output: [42]
}
