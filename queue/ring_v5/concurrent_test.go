package ring_v5

import (
	"fmt"
	gonanoid "github.com/matoous/go-nanoid"
	"reflect"
	"unsafe"
)

import (
	"log"
	"sync"
	"testing"
)

var (
	size = 1000000
)

func Test2P1CKaiXin(t *testing.T) {
	data := generateId()
	buffer := NewRingBuffer[string](1024)

	element := buffer.Element()
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&element))
	dataAddr := unsafe.Pointer(sliceHeader.Data)
	fmt.Printf("address of element: %p\n", dataAddr)

	wg := sync.WaitGroup{}
	go retryProduce(buffer, data)
	wg.Add(1)
	consumer(buffer, data, &wg)
	wg.Wait()
}

func retryProduce(buffer *Ring[string], data map[string]*itemTest) {
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

func consumer(buffer *Ring[string], data map[string]*itemTest, wg *sync.WaitGroup) {
	count := 0
	for {
		id, ok := buffer.Poll()
		if !ok {
			//runtime.Gosched()
			continue
		}
		if id == nil {
			log.Println("nil id type", count)
			break
		}
		count++
		//idStr, ok := (id).(string)
		//if !ok {
		//	log.Println("invalid id type", count, id)
		//	continue
		//	//log.Println("invalid id type", buffer.read, buffer.write)
		//	//continue
		//}
		item, ok := data[*id]
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
		//log.Println("count:", count)
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
