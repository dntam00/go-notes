package pipeline

import (
	"errors"
	"fmt"
	"sync"
)

type Queuer interface {
	Queue(data int) error
	Next() (data, error)
}

type ringBuffer struct {
	scale int
	size  int
	ring  []data
	read  int
	write int
	lock  sync.Mutex
}

type data struct {
	value int
	mask  int
}

func (r *ringBuffer) Queue(value int) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	current := r.write & (r.size - 1)
	d := r.ring[current]
	if d.mask == 1 {
		return errors.New("full")
	}
	r.ring[current] = data{value: value, mask: 1}
	r.write++
	return nil
}

func (r *ringBuffer) Next() (data, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	d := &r.ring[r.read&(r.size-1)]
	if d.mask == 1 {
		r.read += 1
		d.mask = 2
		return *d, nil
	}
	return data{}, errors.New("empty")
}

func (r *ringBuffer) NextSwap() (data, error) {
	//atomic.CompareAndSwapInt32()
	d := &r.ring[r.read&(r.size-1)]
	if d.mask == 1 {
		r.read += 1
		d.mask = 2
		return *d, nil
	}
	return data{}, errors.New("empty")
}

func New(scale int) Queuer {
	size := 1 << scale
	return &ringBuffer{
		scale: scale,
		size:  size,
		ring:  make([]data, size),
	}
}

func main() {
	queuer := New(10)
	_ = queuer.Queue(1)
	_ = queuer.Queue(2)
	_ = queuer.Queue(3)
	_ = queuer.Queue(4)
	fullNums := 0
	for i := 0; i < 1000; i++ {
		err := queuer.Queue(5)
		if err != nil {
			fullNums++
		}
	}
	fmt.Println("error: ", fullNums)
	nextNums := 0
	for i := 0; i < 100; i++ {
		d, err := queuer.Next()
		if err != nil {
			nextNums++
		} else {
			fmt.Println("read: ", d.value)
		}
	}
	fmt.Println("next error: ", nextNums)
}
