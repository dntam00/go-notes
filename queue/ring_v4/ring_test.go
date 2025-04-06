package ring_v4

import (
	"math"
	"testing"
)

func TestRingOverFlowCounter(t *testing.T) {

	type testStructure struct {
		inUse bool
		val   uint64
	}

	capacity := uint64(1024)
	rb := NewRingBuffer(capacity)
	rb.read = math.MaxUint32 - 3
	rb.write = math.MaxUint32 - 3

	size := rb.Size()
	if size != 0 {
		t.Errorf("expected size 0, got %d", rb.Size())
	}

	for i := uint64(0); i < capacity; i++ {
		err := rb.Offer(testStructure{
			inUse: true,
			val:   i,
		})
		if !err {
			t.Errorf("failed to offer: %v", err)
		}
	}

	size = rb.Size()
	if size != capacity {
		t.Errorf("expected size 1024, got %d", rb.Size())
	}

	for i := uint64(0); i < capacity; i++ {
		val, ok := rb.Poll()
		cast, ok := val.(testStructure)
		if !ok {
			t.Errorf("failed to offer")
		} else if cast.val != i {
			t.Errorf("expected val %d, got %d", i, val)
		}
	}
}
