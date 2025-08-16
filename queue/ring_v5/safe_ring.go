package ring_v5

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

type Ring[T any] struct {
	read     uint64
	write    uint64
	element  []*T
	capacity uint64
	mask     uint64
}

func NewRingBuffer[T any](capacity uint64) *Ring[T] {
	return &Ring[T]{
		element:  make([]*T, capacity),
		capacity: capacity,
		mask:     capacity - 1,
		write:    0,
		read:     0,
	}
}

func (r *Ring[T]) Offer(v T) bool {
	oldRead := atomic.LoadUint64(&r.read)
	oldWrite := atomic.LoadUint64(&r.write)
	if r.IsFull(oldRead, oldWrite) {
		//log.Println("full queue")
		return false
	}

	newWrite := oldWrite + 1

	node := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&r.element[newWrite&r.mask])))
	if node != nil {
		return false
	}
	if !atomic.CompareAndSwapUint64(&r.write, oldWrite, newWrite) {
		//log.Println("concurrent write", v, oldRead, oldWrite, newWrite)
		return false
	}

	val := v    // assign to new variable
	ptr := &val // take its address — now val escapes
	//log.Printf("address of stored element: %p\n", ptr)

	//newV := v
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&r.element[newWrite&r.mask])), unsafe.Pointer(ptr))
	return true
}

func (r *Ring[T]) Poll() (*T, bool) {
	oldRead := atomic.LoadUint64(&r.read)
	oldWrite := atomic.LoadUint64(&r.write)
	if r.IsEmpty(oldRead, oldWrite) {
		//log.Println("buffer is empty")
		return nil, false
	}

	newRead := oldRead + 1
	//newReadIndex := newRead & (r.capacity - 1)

	node := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&r.element[newRead&r.mask])))
	if node == nil {
		//log.Println("node is empty")
		return nil, false
	}

	if !atomic.CompareAndSwapUint64(&r.read, oldRead, newRead) {
		//log.Println("concurrent read")
		return nil, false
	}
	//debugPointer(node)
	val := (*T)(node)
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&r.element[newRead&r.mask])), nil)
	return val, true
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

func debugPointer(ptr unsafe.Pointer) {
	if ptr == nil {
		fmt.Println("Nil pointer")
		return
	}

	fmt.Printf("Pointer address: %p\n", ptr)

	casted := (*eface)(ptr)
	if casted == nil {
		fmt.Println("Casted pointer is nil")
		return
	}

	if casted._type.kind != 24 {
		fmt.Printf("Kind: %d\n", casted._type.kind)
		panic("not kind string")
	}

	// Carefully dereference to interface{}
	value := *(*interface{})(ptr)

	if value == nil {
		fmt.Println("Value is nil")
		return
	}
	fmt.Printf("Type: %T\n", value)
	fmt.Printf("Value: %v\n", value)
	fmt.Printf("Size: %v\n", unsafe.Sizeof(value))
}

func (r *Ring[T]) IsFull(read uint64, write uint64) bool {
	return write-read >= r.capacity
}

func (r *Ring[T]) IsEmpty(read uint64, write uint64) bool {
	return (write < read) || (read-write == 0)
}

func (r *Ring[T]) Size() uint64 {
	return r.write - r.read
}

func (r *Ring[T]) Element() []*T {
	return r.element
}
