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
}

func NewRingBuffer(capacity uint64) *RingBuffer {
	return &RingBuffer{
		data:     make([]interface{}, capacity),
		capacity: capacity,
		write:    0,
		read:     0,
	}
}

func (r *RingBuffer) Offer(v interface{}) bool {
	oldRead := atomic.LoadUint64(&r.read)
	oldWrite := atomic.LoadUint64(&r.write)
	if r.isFull(oldRead, oldWrite) {
		return false
	}

	newWrite := oldWrite + 1

	if !atomic.CompareAndSwapUint64(&r.write, oldWrite, newWrite) {
		return false
	}

	newWriteIndex := newWrite & (r.capacity - 1)
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&r.data[newWriteIndex])), unsafe.Pointer(&v))
	return true
}

func (r *RingBuffer) Poll() interface{} {
	oldRead := atomic.LoadUint64(&r.read)
	oldWrite := atomic.LoadUint64(&r.write)
	if r.isEmpty(oldRead, oldWrite) {
		log.Println("buffer is empty")
		return nil
	}

	newRead := oldRead + 1
	if !atomic.CompareAndSwapUint64(&r.read, oldRead, newRead) {
		log.Println("concurrent read")
		return nil
	}
	newReadIndex := newRead & (r.capacity - 1)
	val := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&r.data[newReadIndex])))
	t := *(*interface{})(val)
	return t
}

func (r *RingBuffer) isFull(read uint64, write uint64) bool {
	return write-read >= r.capacity
}

func (r *RingBuffer) isEmpty(read uint64, write uint64) bool {
	return read >= write
}

func (r *RingBuffer) Size() uint64 {
	return r.write - r.read
}
