package main

import (
	"fmt"
	"log"
	"reflect"
	"sync/atomic"
	"unsafe"
)

// https://medium.com/@philbrainy/go-slices-demystified-a-deep-dive-into-memory-layout-and-behavior-59cffd1a49ca
// https://medium.com/codex/go-interface-101-a99943d22bd9

func main() {
	storePtr()
}

type eface struct {
	_type *_type         // Pointer to type information
	data  unsafe.Pointer // Pointer to data
}

type _type struct {
	size       uintptr
	ptrdata    uintptr
	hash       uint32
	tflag      uint8
	align      uint8
	fieldAlign uint8
	kind       uint8
	// These fields help in identifying types
	equal  func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata unsafe.Pointer
}

func storePtr() {
	s := make([]interface{}, 2)

	var v1 interface{}
	v1 = "hello"

	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&s[0])), unsafe.Pointer(&v1))

	log.Printf("address of v1: %p\n", &v1)

	var v2 interface{}
	v2 = "world"
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&s[1])), unsafe.Pointer(&v2))

	log.Printf("address of v1: %p\n", &v2)

	inspectSlice(s)
}

func loadPtr() {
	s := make([]interface{}, 2)
	var v1 interface{}
	v1 = "hello"
	s[0] = v1
	var v2 interface{}
	v2 = "world"
	s[1] = v2

	inspectSlice(s)
}

func inspectSlice(s []interface{}) {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	dataAddr := unsafe.Pointer(sliceHeader.Data)
	log.Printf("address of slice: %p\n", dataAddr)

	log.Println("slice addr:", (unsafe.Pointer)(&s))
	log.Println("slice addr [0]:", (unsafe.Pointer)(&s[0]))
	log.Println("slice addr [1]:", (unsafe.Pointer)(&s[1]))

	//
	loadPtr0 := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s[0])))
	log.Printf("load ptr addr [0]: %p\n", loadPtr0)
	log.Println("load ptr value [0]:", &s[0])
	loadPtr1 := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s[1])))
	log.Printf("load ptr addr [1]: %p\n", loadPtr1)

	//// plus 8 byte to get address of data address
	//valuePtr1 := (unsafe.Pointer)(((uintptr)(unsafe.Pointer(&s[1]))) + 8)
	//val := atomic.LoadPointer((*unsafe.Pointer)(valuePtr1))
	//s1 := *(*string)(val)
	//log.Printf("value slot 1: %s\n", s1)
	//log.Printf("load ptr addr [1]: %p\n", loadPtr1)
	//log.Println("load ptr value [1]:", &s[1])

	dataType := (*eface)(loadPtr0)
	log.Printf("interface{} type: %d\n", dataType._type.kind)
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

func SliceInterface() {
	s := make([]interface{}, 2)
	log.Println("slice addr:", (unsafe.Pointer)(&s))
	fmt.Println("Alignof(interface{}):", unsafe.Alignof(s[0]))
	log.Println("slice [0] addr:", (unsafe.Pointer)(&s[0]))
	log.Println("slice [1] addr:", (unsafe.Pointer)(&s[1]))
	log.Println("space addr:", (uintptr)((unsafe.Pointer)(&s[1]))-(uintptr)((unsafe.Pointer)(&s[0])))
	log.Println("slice converted addr:", (*unsafe.Pointer)((unsafe.Pointer)(&s[0])))

	v1 := "abcxyz"
	var vi1 interface{}
	vi1 = v1
	atomic.StorePointer((*unsafe.Pointer)((unsafe.Pointer)(&s[0])), unsafe.Pointer(&vi1))
	v2 := "hihihaha"
	var vi2 interface{}
	vi2 = v2
	atomic.StorePointer((*unsafe.Pointer)((unsafe.Pointer)(&s[1])), unsafe.Pointer(&vi2))

	val := *(*interface{})(atomic.LoadPointer((*unsafe.Pointer)((unsafe.Pointer)(&s[0]))))
	log.Println("slice [0] value:", val)

	atomic.StorePointer((*unsafe.Pointer)((unsafe.Pointer)(&s[1])), unsafe.Pointer(nil))

	pointer := atomic.LoadPointer((*unsafe.Pointer)((unsafe.Pointer)(&s[1])))
	if pointer == nil {
		log.Println("slice [1] is nil")
		return
	}
	val = *(*interface{})(pointer)
	log.Println("slice [1] value:", val)
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
