package ring_v4

import (
	"fmt"
	"reflect"
	"sync/atomic"
	"unsafe"
)

// https://www.lenshood.dev/2021/04/19/lock-free-ring-buffer/

type Ring struct {
	read  uint64
	write uint64
	//element  []unsafe.Pointer
	element  []interface{}
	capacity uint64
	mask     uint64
}

func NewRingBuffer(capacity uint64) *Ring {
	return &Ring{
		//element:  make([]unsafe.Pointer, capacity),
		element:  make([]interface{}, capacity),
		capacity: capacity,
		mask:     capacity - 1,
		write:    0,
		read:     0,
	}
}

func (r *Ring) Offer(v interface{}) bool {
	oldRead := atomic.LoadUint64(&r.read)
	oldWrite := atomic.LoadUint64(&r.write)
	if r.isFull(oldRead, oldWrite) {
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

	//newV := v
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&r.element[newWrite&r.mask])), unsafe.Pointer(ptr))
	return true
}

func (r *Ring) Poll() (interface{}, bool) {
	oldRead := atomic.LoadUint64(&r.read)
	oldWrite := atomic.LoadUint64(&r.write)
	if r.isEmpty(oldRead, oldWrite) {
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
	debugPointer(node)
	val := *(*interface{})(node)
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&r.element[newRead&r.mask])), nil)
	return val, true
}

func debugPointer(ptr unsafe.Pointer) {
	if ptr == nil {
		fmt.Println("Nil pointer")
		return
	}

	// Carefully dereference to interface{}
	value := *(*interface{})(ptr)

	fmt.Printf("Pointer address: %p\n", ptr)
	if value == nil {
		fmt.Println("Value is nil")
		return
	}
	fmt.Printf("Type: %T\n", value)
	fmt.Printf("Value: %v\n", value)

	// You can also use reflection for more details
	valueType := reflect.TypeOf(value)
	fmt.Printf("Reflection type: %v\n", valueType)
	fmt.Printf("Kind: %v\n", valueType.Kind())
}

func (r *Ring) isFull(read uint64, write uint64) bool {
	return write-read >= r.capacity
}

func (r *Ring) isEmpty(read uint64, write uint64) bool {
	return (write < read) || (read-write == 0)
}

func (r *Ring) Size() uint64 {
	return r.write - r.read
}
