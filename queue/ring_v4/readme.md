```
2025/04/06 15:43:21 concurrent write wIIwY7v3WcLMTJB64--Be 395830 396853 396854
==================
WARNING: DATA RACE
Read at 0x00c010073a40 by goroutine 6:
  play-around/queue/ring_v4.(*RingBuffer).Poll()
      /Users/dntam/Projects/golang/go-notes/queue/ring_v4/ring.go:73 +0x170
  play-around/queue/ring_v4.consumer()
      /Users/dntam/Projects/golang/go-notes/queue/ring_v4/concurrent_test.go:81 +0x44
  play-around/queue/ring_v4.Test2P1C()
      /Users/dntam/Projects/golang/go-notes/queue/ring_v4/concurrent_test.go:20 +0x1e0
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1690 +0x1a8
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1743 +0x40

Previous write at 0x00c010073a40 by goroutine 2464:
  play-around/queue/ring_v4.(*RingBuffer).Offer()
      /Users/dntam/Projects/golang/go-notes/queue/ring_v4/ring.go:27 +0x44
  play-around/queue/ring_v4.produce.func1()
      /Users/dntam/Projects/golang/go-notes/queue/ring_v4/concurrent_test.go:58 +0x54
  play-around/queue/ring_v4.produce.gowrap1()
      /Users/dntam/Projects/golang/go-notes/queue/ring_v4/concurrent_test.go:66 +0x54

Goroutine 6 (running) created at:
  testing.(*T).Run()
      /usr/local/go/src/testing/testing.go:1743 +0x674
  testing.runTests.func1()
      /usr/local/go/src/testing/testing.go:2168 +0x80
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1690 +0x1a8
  testing.runTests()
      /usr/local/go/src/testing/testing.go:2166 +0x764
  testing.(*M).Run()
      /usr/local/go/src/testing/testing.go:2034 +0xba0
  main.main()
      _testmain.go:49 +0x110

Goroutine 2464 (running) created at:
  play-around/queue/ring_v4.produce()
      /Users/dntam/Projects/golang/go-notes/queue/ring_v4/concurrent_test.go:57 +0x2e4
  play-around/queue/ring_v4.Test2P1C.gowrap1()
      /Users/dntam/Projects/golang/go-notes/queue/ring_v4/concurrent_test.go:18 +0x40
==================
runtime: pointer 0xc00fdb5530 to unallocated span span.base()=0xc00fdb4000 span.limit=0xc00fdb6000 span.state=0
fatal error: found bad pointer in Go heap (incorrect use of unsafe or cgo?)
```