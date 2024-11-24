package main

import (
	"math/rand"
	"sync"
	"sync/atomic"
)

func main() {
	nums := make([]int, 10000000)
	sumSequential(nums)
	sumConcurrent(8, nums)
}

func generateList(totalNumbers int) []int {
	numbers := make([]int, totalNumbers)
	for i := 0; i < totalNumbers; i++ {
		numbers[i] = rand.Intn(totalNumbers)
	}
	return numbers
}

func sumConcurrent(goroutine int, nums []int) int64 {
	partition := len(nums) / goroutine
	var wg sync.WaitGroup
	wg.Add(goroutine)
	//start := time.Now()
	var v int64
	for i := 0; i < goroutine; i++ {
		go func(index int) {
			sum := 0
			start := index * partition
			end := start + partition
			if index == goroutine-1 {
				end = len(nums)
			}
			for _, n := range nums[start:end] {
				sum += n
			}
			atomic.AddInt64(&v, int64(sum))
			wg.Done()
		}(i)
	}
	wg.Wait()
	return v
	//fmt.Println("concurrent take: ", time.Since(start))
}

func sumSequential(nums []int) int {
	//start := time.Now()
	var v int
	for _, n := range nums {
		v += n
	}
	return v
	//fmt.Println("sequential take: ", time.Since(start))
}
