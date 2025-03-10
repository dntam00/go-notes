package cpu_cache

import (
	"sync/atomic"
	"testing"
)

func BenchmarkPadding(b *testing.B) {
	b.Run("without_padding", func(b *testing.B) {
		nums := [128]atomic.Int64{}
		i := atomic.Int64{}
		b.RunParallel(func(pb *testing.PB) {
			id := i.Add(1)
			for pb.Next() {
				nums[id].Add(1)
			}
		})
	})
	b.Run("with_padding", func(b *testing.B) {
		type pad struct {
			val atomic.Int64
			_   [8]uint64
		}
		nums := [128]pad{}
		i := atomic.Int64{}
		b.RunParallel(func(pb *testing.PB) {
			id := i.Add(1)
			for pb.Next() {
				nums[id].val.Add(1)
			}
		})
	})
}
