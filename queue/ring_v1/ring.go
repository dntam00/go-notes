package ring_v1

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
	return r.read == r.capacity
}

func (r *RingBuffer[T]) Capacity() int {
	return r.capacity
}

func (r *RingBuffer[T]) Size() int {
	return r.write - r.read
}

func (r *RingBuffer[T]) Offer(item T) error {
	if r.write-r.read >= r.capacity {
		return ErrQueueFull
	}
	r.data[r.write%r.capacity] = item
	r.write++
	return nil
}

func (r *RingBuffer[T]) Poll() (*T, error) {
	if r.read == r.write {
		return nil, ErrQueueEmpty
	}
	t := r.data[r.read%r.capacity]
	r.read++
	return &t, nil
}
