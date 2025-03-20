package main

import (
	"errors"
	"fmt"
	"log"
	"syscall"
	"time"
)

const (
	MaxEvents      = 10
	ReadBufferSize = 1
)

func main() {
	// Create epoll instance
	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(epfd)

	// Create a pipe
	var pipeFds [2]int
	if err := syscall.Pipe(pipeFds[:]); err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(pipeFds[0])
	defer syscall.Close(pipeFds[1])

	// Add read end of the pipe to epoll with edge-triggered mode
	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN | syscall.EPOLLET,
		Fd:     int32(pipeFds[0]),
	}
	if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, pipeFds[0], &event); err != nil {
		log.Fatal(err)
	}

	// Write to the pipe in a separate goroutine
	go func() {
		//for {
		time.Sleep(1 * time.Second) // Simulate a delay
		data := []byte("Hello, epoll!")
		if _, err := syscall.Write(pipeFds[1], data); err != nil {
			log.Println("Error writing to pipe:", err)
		}
		//}
	}()

	// Create buffer for events
	events := make([]syscall.EpollEvent, MaxEvents)

	// Event loop
	for {
		n, err := syscall.EpollWait(epfd, events, -1)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("epoll returned with %d readble file descriptior\n", n)

		// comment to read data
		// this block of code show that level-triggered return true if the file descriptor is still readable
		//if n != 0 {
		//	log.Printf("epoll returned with %d readble file descriptior\n", n)
		//	time.Sleep(1 * time.Second)
		//	continue
		//}

		for i := 0; i < n; i++ {
			if events[i].Events&syscall.EPOLLIN != 0 {
				buffer := make([]byte, ReadBufferSize)
				//for {
				n, err := syscall.Read(int(events[i].Fd), buffer)
				if err != nil {
					if errors.Is(err, syscall.EAGAIN) {
						break
					}
					log.Println("Error reading from pipe:", err)
					break
				}
				if n == 0 {
					break
				}
				fmt.Printf("Read %d bytes: %s\n", n, string(buffer[:n]))
				//}
			}
		}
	}
}
