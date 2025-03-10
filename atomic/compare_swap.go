package main

import (
	"runtime"
	"sync/atomic"
	"time"
)

type request struct {
	mask uint32
}

func WaitCompareSwap() int64 {
	r := &request{}
	r.mask = 1
	now := time.Now()
	ch := make(chan struct{})

	go func() {
		for !atomic.CompareAndSwapUint32(&r.mask, 0, 1) {
			runtime.Gosched()
		}
		ch <- struct{}{}
	}()

	time.Sleep(3 * time.Second)
	atomic.CompareAndSwapUint32(&r.mask, 1, 0)
	<-ch
	return time.Since(now).Milliseconds()
}
