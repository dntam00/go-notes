package ring_v2

import "errors"

var (
	ErrQueueFull  = errors.New("queue is full")
	ErrQueueEmpty = errors.New("queue is empty")
)

type RingBuffer[T any] struct {
	data     []T
	read     uint32
	write    uint32
	capacity uint32
}

func NewRingBuffer[T any](capacity uint32) *RingBuffer[T] {
	return &RingBuffer[T]{
		data:     make([]T, capacity),
		capacity: capacity,
		write:    0,
		read:     0,
	}
}

func (r *RingBuffer[T]) Mask(v uint32) uint32 {
	return v & (r.Capacity() - 1)
}

func (r *RingBuffer[T]) IsEmpty() bool {
	return r.read == r.write
}

func (r *RingBuffer[T]) IsFull() bool {
	return r.Size() == r.Capacity()
}

func (r *RingBuffer[T]) Capacity() uint32 {
	return r.capacity
}

func (r *RingBuffer[T]) Size() uint32 {
	return r.write - r.read
}

func (r *RingBuffer[T]) Offer(item T) error {
	if r.IsFull() {
		return ErrQueueFull
	}
	r.data[r.Mask(r.write)] = item
	r.write = r.write + 1
	return nil
}

func (r *RingBuffer[T]) Poll() (*T, error) {
	if r.IsEmpty() {
		return nil, ErrQueueEmpty
	}
	t := r.data[r.Mask(r.read)]
	r.read = r.read + 1
	return &t, nil
}
