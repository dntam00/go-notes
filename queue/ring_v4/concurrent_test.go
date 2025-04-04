package ring_v4

import (
	gonanoid "github.com/matoous/go-nanoid"
	"sync"
	"testing"
)

var (
	size = 100000
)

func Test2P1C(*testing.T) {
	data := generateId()
	buffer := NewRingBuffer(1024)
	wg := sync.WaitGroup{}
	go produce(buffer, data)
	wg.Add(1)
	consumer(buffer, data, &wg)
	wg.Wait()
}

func produce(buffer *RingBuffer, data map[string]*itemTest) {
	for _, v := range data {
		buffer.Offer(v.id)
	}
}

func consumer(buffer *RingBuffer, data map[string]*itemTest, wg *sync.WaitGroup) {
	count := 0
	for {
		id := buffer.Poll()
		if id == nil {
			continue
		}
		count++
		idStr := id.(string)
		item, ok := data[idStr]
		if !ok {
			panic("item not found")
		}
		if item.read >= 1 {
			panic("expect read less than 1")
		}
		item.read = 1
		if count == size {
			break
		}
	}
	wg.Done()
}

type itemTest struct {
	id    string
	write int32
	read  int32
}

func generateId() map[string]*itemTest {
	data := make(map[string]*itemTest)
	for i := 0; i < size; i++ {
		id, err := gonanoid.Nanoid()
		if err != nil {
			panic(err)
		}
		data[id] = &itemTest{
			id:    id,
			write: 0,
			read:  0,
		}
	}
	return data
}
