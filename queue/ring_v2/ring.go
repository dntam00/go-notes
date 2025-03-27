package ring_v2

import "errors"

var (
	ErrQueueFull  = errors.New("queue is full")
	ErrQueueEmpty = errors.New("queue is empty")
)

type RingBuffer[T any] struct {
	data     []T
	read     int
	write    int
	capacity int
}

func NewRingBuffer[T any](size int) *RingBuffer[T] {
	return &RingBuffer[T]{
		data:     make([]T, size),
		capacity: size,
		write:    0,
		read:     0,
	}
}

func (r *RingBuffer[T]) IsEmpty() bool {
	return r.read == r.write
}

func (r *RingBuffer[T]) IsFull() bool {
	return r.Size() == r.Capacity()
}

func (r *RingBuffer[T]) Capacity() int {
	return r.capacity - 1
}

func (r *RingBuffer[T]) Size() int {
	return (r.write - r.read) & r.Capacity()
}

func (r *RingBuffer[T]) Offer(item T) error {
	if r.IsFull() {
		return ErrQueueFull
	}
	r.data[r.write] = item
	r.write = (r.write + 1) & r.Capacity()
	return nil
}

func (r *RingBuffer[T]) Poll() (*T, error) {
	if r.IsEmpty() {
		return nil, ErrQueueEmpty
	}
	t := r.data[r.read]
	r.read = (r.read + 1) & r.Capacity()
	return &t, nil
}
