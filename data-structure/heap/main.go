package main

import (
	"errors"
	"fmt"
	"math"
)

type heap struct {
	arr []int
}

func (h *heap) upHeap() {
	index := len(h.arr) - 1
	for index != 0 {
		p := (index - 1) / 2
		if h.arr[index] >= h.arr[p] {
			break
		}
		t := h.arr[index]
		h.arr[index] = h.arr[p]
		h.arr[p] = t
		index = p
	}
}

func (h *heap) downHeap() {
	c := 0
	size := len(h.arr)
	for {
		// 0
		left := (c+1)*2 - 1
		right := (c + 1) * 2

		if left >= size {
			return
		}

		swapIndex := left

		if right < size {
			if h.arr[right] < h.arr[left] {
				swapIndex = right
			}
		}

		currentVal := h.arr[c]
		if currentVal > h.arr[swapIndex] {
			// swap
			h.arr[c] = h.arr[swapIndex]
			h.arr[swapIndex] = currentVal
			c = swapIndex
		} else if currentVal <= h.arr[swapIndex] {
			// less than its children
			return
		}
	}
}

func (h *heap) poll() (int, error) {
	if len(h.arr) == 0 {
		return 0, errors.New("empty")
	}
	returned := h.arr[0]
	h.arr[0] = h.arr[len(h.arr)-1]
	h.arr = h.arr[:len(h.arr)-1]
	h.downHeap()
	return returned, nil
}

func (h *heap) peak() (int, error) {
	if len(h.arr) == 0 {
		return 0, errors.New("empty")
	}
	return h.arr[0], nil
}

func (h *heap) offer(val int) {
	h.arr = append(h.arr, val)
	h.upHeap()
}

func minOperations(nums []int, k int) int {
	h := heap{
		arr: make([]int, 0),
	}
	for _, c := range nums {
		h.offer(c)
	}
	operation := 0
	for {
		x, _ := h.poll()
		if x >= k {
			return operation
		}
		operation++
		y, _ := h.poll()
		//fmt.Println(x, y)
		val := math.Max(float64(x), float64(y)) + math.Min(float64(x), float64(y))*2
		h.offer(int(val))
	}
}

func main() {
	arr := []int{1, 2, 3}
	fmt.Println(minOperations(arr, 10))
	//test(arr)
}

func test(arr []int) {
	h := heap{}
	h.arr = arr
	h.upHeap()

	for i := 0; i < len(arr); i++ {
		polled, err := h.poll()
		if err != nil {
			break
		}
		fmt.Println(polled)
	}
}
