package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func BenchmarkSequential(b *testing.B) {
	for i := 0; i < 1; i++ {
		queuer := New(15)
		successQueue := 0
		failNext := 0
		wg := sync.WaitGroup{}
		wg.Add(2)
		start := time.Now()
		go func() {
			for i := 0; i < 100000; i++ {
				time.Sleep(5 * time.Millisecond)
				err := queuer.Queue(i)
				if err == nil {
					successQueue++
				}
			}
			wg.Done()
		}()
		time.Sleep(500 * time.Millisecond)
		l := make([]int, 100000)
		index := 0
		go func() {
			pre := -1
			for i := 0; i < 100000; i++ {
				time.Sleep(5 * time.Millisecond)
				v, err := queuer.Next()
				if err == nil {
					value := v.value
					if pre > value {
						//fmt.Println("failed: pre: ", pre, "v: ", v.value)
						break
					}
					//fmt.Println("success: pre before: ", pre, "v: ", v.value)
					l[index] = v.value
					index++
					pre = v.value
					//fmt.Println("success: pre after: ", pre, "v: ", value, v.value)
				} else {
					failNext++
				}
			}
			wg.Done()
		}()
		wg.Wait()
		fmt.Println("need time: ", time.Since(start).Seconds())
		pre := -1
		check := true
		fmt.Println("fail get: ", failNext)
		fmt.Println("success queue: ", successQueue)
		for i, d := range l {
			if i < index {
				if pre > d {
					fmt.Println("pre: ", pre, "d: ", d)
					check = false
					break
				}
				pre = d
			} else {
				break
			}
		}
		fmt.Println("size: ", index)
		assert.True(b, check, "consume message correctly")
	}
}
