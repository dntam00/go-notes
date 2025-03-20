package main

import (
	"fmt"
	"unsafe"
)

// https://medium.com/@bradford_hamilton/exploring-unsafe-features-in-go-1-20-a-hands-on-demo-7149ba82e6e1

func main() {
	// Allocate integer, get pointer to it:
	x := 10
	xPtr := &x

	// Get a uintptr of the address of x. Do NOT do this.
	xUintPtr := uintptr(unsafe.Pointer(xPtr))

	// ---------------------------------------------------------------
	// At this point, `x` is unused and so could be garbage collected.
	// If that happens, we then have an uintptr (integer) that when
	// casted back to an unsafe.Pointer, points to to some invalid
	// piece of memory.
	// ---------------------------------------------------------------

	fmt.Println(*(*int)(unsafe.Pointer(xUintPtr))) // possible misuse of unsafe.Pointer
}
