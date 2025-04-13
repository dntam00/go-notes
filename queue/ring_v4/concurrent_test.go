package ring_v4

import (
	gonanoid "github.com/matoous/go-nanoid"
)

import (
	"log"
	"sync"
	"testing"
)

var (
	size = 10000000
)

func Test1P1C(b *testing.T) {
	for i := 0; i < 1; i++ {
		data := generateId()
		buffer := NewRingBuffer(1024)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go singleProduce(buffer, data)
		consumer(buffer, data, &wg)
		wg.Wait()
	}
}

func Test2P1CKaiXin(t *testing.T) {
	data := generateId()
	buffer := NewRingBuffer(1024)
	wg := sync.WaitGroup{}
	go retryProduce(buffer, data)
	wg.Add(1)
	consumer(buffer, data, &wg)
	wg.Wait()
}

func singleProduce(buffer *Ring, data map[string]*itemTest) {
	for _, v := range data {
		id := v.id
		offer := false
		for !offer {
			offer = buffer.Offer(id)
		}
	}
	log.Println("finish retryProduce")
}

func retryProduce(buffer *Ring, data map[string]*itemTest) {
	for _, v := range data {
		id := v.id
		offer := buffer.Offer(id)
		if !offer {
			go func(idStr string) {
				for !buffer.Offer(idStr) {

				}
				item, ok := data[idStr]
				if !ok {
					panic("idStr not found")
				}
				item.write++
			}(id)
			continue
		}
		item, ok := data[id]
		if !ok {
			panic("item not found")
		}
		item.write++
	}
	log.Println("finish retryProduce")
}

func consumer(buffer *Ring, data map[string]*itemTest, wg *sync.WaitGroup) {
	count := 0
	for {
		id, ok := buffer.Poll()
		if !ok {
			//runtime.Gosched()
			continue
		}
		count++
		idStr, ok := id.(string)
		if id == nil {
			log.Println("nil id type", count)
			break
		}
		if !ok {
			log.Println("invalid id type", count, id)
			continue
			//log.Println("invalid id type", buffer.read, buffer.write)
			//continue
		}
		item, ok := data[idStr]
		if !ok {
			panic("item not found")
		}
		if item.read >= 1 {
			log.Println("item written 1:", item)
			//time.Sleep(time.Millisecond * 5000)
			//log.Println("item written 2:", item)
			panic("expect read less than 1")
		}
		item.read++
		if count == size {
			break
		}
	}
	log.Println("finish consumer")
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
