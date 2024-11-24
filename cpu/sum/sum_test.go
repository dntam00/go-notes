package main

import (
	"runtime"
	"testing"
)

func BenchmarkSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sumSequential(generateList(10000000))
	}
}

func BenchmarkConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sumConcurrent(runtime.NumCPU(), generateList(10000000))
	}
}
