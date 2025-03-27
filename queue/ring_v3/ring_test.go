package ring_v2

import (
	"errors"
	"testing"
)

func TestRing(t *testing.T) {
	size := uint32(4)
	rb := NewRingBuffer[string](size)

	if rb.Capacity() != size {
		t.Errorf("Expected capacity %d, got %d", size, rb.Capacity())
	}

	if !rb.IsEmpty() {
		t.Error("New buffer should be empty")
	}

	if rb.Size() != 0 {
		t.Errorf("Expected size 0, got %d", rb.Size())
	}

	// Test offering items
	err := rb.Offer("item1")
	if err != nil {
		t.Errorf("Failed to offer item: %v", err)
	}

	if rb.Size() != 1 {
		t.Errorf("Expected size 1, got %d", rb.Size())
	}

	// Test polling items
	item, err := rb.Poll()
	if err != nil {
		t.Errorf("Failed to poll item: %v", err)
	}

	if *item != "item1" {
		t.Errorf("Expected 'item1', got '%s'", *item)
	}

	// Test polling from empty buffer
	_, err = rb.Poll()
	if !errors.Is(err, ErrQueueEmpty) {
		t.Errorf("Expected ErrQueueEmpty, got %v", err)
	}

	// Test filling the buffer
	for i := uint32(0); i < size; i++ {
		rb.Offer("item")
	}

	// Test offering to a full buffer
	err = rb.Offer("overflow")
	if !errors.Is(err, ErrQueueFull) {
		t.Errorf("Expected ErrQueueFull, got %v", err)
	}
}

func TestRingBuffer(t *testing.T) {
	// Test FIFO behavior
	rb := NewRingBuffer[int](4)
	rb.Offer(1)
	rb.Offer(2)
	rb.Offer(3)
	rb.Offer(4)

	val1, _ := rb.Poll()
	val2, _ := rb.Poll()
	val3, _ := rb.Poll()
	val4, _ := rb.Poll()

	if *val1 != 1 || *val2 != 2 || *val3 != 3 || *val4 != 4 {
		t.Error("Buffer does not maintain FIFO order")
	}

	_, err := rb.Poll()
	if !errors.Is(err, ErrQueueEmpty) {
		t.Errorf("Expected ErrQueueEmpty, got %v", err)
	}
}
