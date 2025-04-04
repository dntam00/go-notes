package main

import (
	"log"
	"math"
	"sync/atomic"
	"unsafe"
)

func main() {

	//log.Println((u - 3) & 1023)
	//log.Println(u & 1023)
	//log.Println((u + 1) & 1023)
	//log.Println(u + 2 - u)
	//log.Println(u + 1)

	//log.Println(findPowerOfTwo(1022))
	//log.Println(strconv.FormatInt(42, 2))
	//unsafePointerUsage()
	//SliceAddr()
	overflow()
}

func overflow() {
	u := uint32(0)
	u = math.MaxUint32
	log.Println(0 - u)
}

type ring struct {
	element []interface{}
}

func SliceAddr() {
	ints := make([]int, 2)
	log.Println("ints addr:", (unsafe.Pointer)(&ints))
	log.Println("ints addr 0:", (unsafe.Pointer)(&ints[0]))
	log.Println("ints addr 1:", (unsafe.Pointer)(&ints[1]))
	log.Println("ints addr converted:", (*unsafe.Pointer)((unsafe.Pointer)(&ints[0])))
}

func unsafePointerUsage() {
	var v interface{}
	v = "abcxyz"
	r := ring{element: make([]interface{}, 1)}
	ptr := (*unsafe.Pointer)(unsafe.Pointer(&r.element[0]))
	atomic.StorePointer(ptr, unsafe.Pointer(&v))
	val := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&r.element[0])))
	log.Println(r.element[0])
	log.Println(*(*interface{})(val))
}

func findPowerOfTwo(givenMum uint64) uint64 {
	givenMum--
	givenMum |= givenMum >> 1
	givenMum |= givenMum >> 2
	givenMum |= givenMum >> 4
	givenMum |= givenMum >> 8
	givenMum |= givenMum >> 16
	givenMum |= givenMum >> 32
	givenMum++

	return givenMum
}
