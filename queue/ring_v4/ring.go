package ring_v4

import (
	"log"
	"sync/atomic"
	"unsafe"
)

// https://www.lenshood.dev/2021/04/19/lock-free-ring-buffer/

type RingBuffer struct {
	read     uint64
	write    uint64
	data     []interface{}
	capacity uint64
	mask     uint64
}

func NewRingBuffer(capacity uint64) *RingBuffer {
	return &RingBuffer{
		data:     make([]interface{}, capacity),
		capacity: capacity,
		mask:     capacity - 1,
		write:    0,
		read:     0,
	}
}

func (r *RingBuffer) Offer(v interface{}) bool {
	oldRead := atomic.LoadUint64(&r.read)
	oldWrite := atomic.LoadUint64(&r.write)
	if r.isFull(oldRead, oldWrite) {
		//log.Println("full queue")
		return false
	}

	newWrite := oldWrite + 1

	node := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&r.data[newWrite&r.mask])))
	if node != nil {
		return false
	}
	if !atomic.CompareAndSwapUint64(&r.write, oldWrite, newWrite) {
		log.Println("concurrent write", v, oldRead, oldWrite, newWrite)
		return false
	}

	val := v    // assign to new variable
	ptr := &val // take its address — now val escapes

	//newV := v
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&r.data[newWrite&r.mask])), unsafe.Pointer(ptr))
	return true
}

func (r *RingBuffer) Poll() (interface{}, bool) {
	oldRead := atomic.LoadUint64(&r.read)
	oldWrite := atomic.LoadUint64(&r.write)
	if r.isEmpty(oldRead, oldWrite) {
		log.Println("buffer is empty")
		return nil, false
	}

	newRead := oldRead + 1
	//newReadIndex := newRead & (r.capacity - 1)

	node := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&r.data[newRead&r.mask])))
	if node == nil {
		log.Println("node is empty")
		return nil, false
	}

	if !atomic.CompareAndSwapUint64(&r.read, oldRead, newRead) {
		log.Println("concurrent read")
		return nil, false
	}
	val := *(*interface{})(node)
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&r.data[newRead&r.mask])), nil)
	return val, true
}

func (r *RingBuffer) isFull(read uint64, write uint64) bool {
	return write-read >= r.capacity
}

func (r *RingBuffer) isEmpty(read uint64, write uint64) bool {
	return (write < read) || (read-write == 0)
}

func (r *RingBuffer) Size() uint64 {
	return r.write - r.read
}
