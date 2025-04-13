package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

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
	equal     func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata    unsafe.Pointer
	str       nameOff
	ptrToThis typeOff
}

type nameOff int32
type typeOff int32

type interfaceHeader struct {
	typ  *_type
	data unsafe.Pointer
}

const (
	tflagUncommon      = 1 << 0
	tflagExtraStar     = 1 << 1
	tflagNamed         = 1 << 2
	tflagRegularMemory = 1 << 3
)

const (
	kindBool = 1 + iota
	kindInt
	kindInt8
	kindInt16
	kindInt32
	kindInt64
	kindUint
	kindUint8
	kindUint16
	kindUint32
	kindUint64
	kindUintptr
	kindFloat32
	kindFloat64
	kindComplex64
	kindComplex128
	kindArray
	kindChan
	kindFunc
	kindInterface
	kindMap
	kindPtr
	kindSlice
	kindString
	kindStruct
	kindUnsafePointer
)

// Maps runtime kind to string representation
var kindNames = [...]string{
	kindBool:          "bool",
	kindInt:           "int",
	kindInt8:          "int8",
	kindInt16:         "int16",
	kindInt32:         "int32",
	kindInt64:         "int64",
	kindUint:          "uint",
	kindUint8:         "uint8",
	kindUint16:        "uint16",
	kindUint32:        "uint32",
	kindUint64:        "uint64",
	kindUintptr:       "uintptr",
	kindFloat32:       "float32",
	kindFloat64:       "float64",
	kindComplex64:     "complex64",
	kindComplex128:    "complex128",
	kindArray:         "array",
	kindChan:          "chan",
	kindFunc:          "func",
	kindInterface:     "interface",
	kindMap:           "map",
	kindPtr:           "ptr",
	kindSlice:         "slice",
	kindString:        "string",
	kindStruct:        "struct",
	kindUnsafePointer: "unsafe.Pointer",
}

func main() {
	// Create a slice of interface{} with different types
	s := make([]interface{}, 6)
	s[0] = 42                            // int
	s[1] = "hello"                       // string
	s[2] = 3.14                          // float64
	s[3] = true                          // bool
	s[4] = []int{1, 2, 3}                // slice
	s[5] = struct{ name string }{"test"} // struct

	fmt.Println("Slice details:")
	fmt.Printf("Slice header address: %p\n", &s)
	fmt.Printf("Size of interface{}: %d bytes\n", unsafe.Sizeof(s[0]))
	fmt.Println()

	// Inspect each element in the slice
	for i := range s {
		fmt.Printf("Element %d (%v):\n", i, s[i])

		// Get the interface header using unsafe
		ifaceHeader := (*interfaceHeader)(unsafe.Pointer(&s[i]))
		typeInfo := ifaceHeader.typ

		// Print memory addresses
		fmt.Printf("  Interface header address: %p\n", &s[i])
		fmt.Printf("  Type pointer: %p\n", typeInfo)
		fmt.Printf("  Data pointer: %p\n", ifaceHeader.data)

		// Decode type flags
		fmt.Printf("  Type size: %d bytes\n", typeInfo.size)
		fmt.Printf("  Type hash: 0x%x\n", typeInfo.hash)

		// Decode type flags
		fmt.Printf("  Type flags: 0x%x (", typeInfo.tflag)
		if typeInfo.tflag&tflagUncommon != 0 {
			fmt.Print("Uncommon ")
		}
		if typeInfo.tflag&tflagExtraStar != 0 {
			fmt.Print("ExtraStar ")
		}
		if typeInfo.tflag&tflagNamed != 0 {
			fmt.Print("Named ")
		}
		if typeInfo.tflag&tflagRegularMemory != 0 {
			fmt.Print("RegularMemory")
		}
		fmt.Println(")")

		// Decode kind
		fmt.Printf("  Kind value: %d\n", typeInfo.kind)
		if int(typeInfo.kind) < len(kindNames) && kindNames[typeInfo.kind] != "" {
			fmt.Printf("  Kind name: %s\n", kindNames[typeInfo.kind])
		} else {
			fmt.Printf("  Kind name: unknown(%d)\n", typeInfo.kind)
		}

		// Use reflect to get type information safely
		elemType := reflect.TypeOf(s[i])
		elemValue := reflect.ValueOf(s[i])

		fmt.Printf("  Reflect Type: %s\n", elemType.String())
		fmt.Printf("  Reflect Kind: %s\n", elemType.Kind())
		fmt.Printf("  Reflect Size: %d bytes\n", elemType.Size())

		// Additional type information
		fmt.Printf("  Alignment: %d bytes\n", typeInfo.align)
		fmt.Printf("  Field alignment: %d bytes\n", typeInfo.fieldAlign)
		fmt.Printf("  Pointer data size: %d bytes\n", typeInfo.ptrdata)

		// Type-specific information
		switch elemType.Kind() {
		case reflect.Struct:
			fmt.Printf("  Number of fields: %d\n", elemType.NumField())
			for j := 0; j < elemType.NumField(); j++ {
				field := elemType.Field(j)
				fmt.Printf("    Field %d: %s %s (offset: %d)\n",
					j, field.Name, field.Type, field.Offset)
			}
		case reflect.Slice:
			fmt.Printf("  Element type: %s\n", elemType.Elem())
			fmt.Printf("  Length: %d\n", elemValue.Len())
			fmt.Printf("  Capacity: %d\n", elemValue.Cap())
		}

		fmt.Println()
	}

	// Demonstrate memory layout
	fmt.Println("Memory layout of the slice:")
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	fmt.Printf("  Slice data pointer: %p\n", unsafe.Pointer(sliceHeader.Data))
	fmt.Printf("  Slice length: %d\n", sliceHeader.Len)
	fmt.Printf("  Slice capacity: %d\n", sliceHeader.Cap)

	// Accessing the backing array
	fmt.Println("\nExamining backing array elements:")
	for i := 0; i < sliceHeader.Len; i++ {
		// Calculate the address of each element in the backing array
		elemAddr := unsafe.Pointer(sliceHeader.Data + uintptr(i)*unsafe.Sizeof(s[0]))
		elemHeader := (*interfaceHeader)(elemAddr)
		fmt.Printf("  Element %d address: %p\n", i, elemAddr)
		fmt.Printf("  Element %d type pointer: %p\n", i, elemHeader.typ)
		fmt.Printf("  Element %d data pointer: %p\n", i, elemHeader.data)
	}
}
