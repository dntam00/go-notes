package main

import (
	"fmt"
	"syscall"
	"time"
)

// https://www.sobyte.net/post/2022-01/go-netpoller/#io-multiplexing
// https://www.sobyte.net/post/2022-05/linux-blocking-io/
// https://man7.org/linux/man-pages/man2/select.2.html
// https://jvns.ca/blog/2017/06/03/async-io-on-linux--select--poll--and-epoll/

func main() {
	// Create a pipe
	var pipeFds [2]int
	err := syscall.Pipe(pipeFds[:])
	if err != nil {
		fmt.Println("Error creating pipe:", err)
		return
	}
	defer syscall.Close(pipeFds[0])
	defer syscall.Close(pipeFds[1])

	// Write some data to the pipe
	go func() {
		time.Sleep(1 * time.Second) // Simulate a delay
		data := []byte("Hello, world!")
		_, err := syscall.Write(pipeFds[1], data)
		if err != nil {
			fmt.Println("Error writing to pipe:", err)
			return
		}
		fmt.Println("Data written to pipe.")
	}()

	// Monitor the read end of the pipe using syscall.Select
	for {
		// Set up the file descriptor set for select
		var readFds syscall.FdSet
		readFds.Bits[pipeFds[0]/64] |= 1 << (pipeFds[0] % 64)

		// Call select with a timeout (5 seconds)
		timeout := syscall.Timeval{Sec: 5, Usec: 0}
		n, err := syscall.Select(pipeFds[0]+1, &readFds, nil, nil, &timeout)
		if err != nil {
			fmt.Println("Error in select:", err)
			return
		}

		if n == 0 {
			fmt.Println("Select timed out.")
			continue
		}

		// Check if the pipe is readable
		if readFds.Bits[pipeFds[0]/64]&(1<<(pipeFds[0]%64)) != 0 {
			fmt.Println("Pipe is readable.")

			// Uncomment the following code to read the data and avoid busy loop
			/*
				buf := make([]byte, 1024)
				n, err := syscall.Read(pipeFds[0], buf)
				if err != nil {
					fmt.Println("Error reading from pipe:", err)
					return
				}
				fmt.Printf("Read %d bytes: %s\n", n, string(buf[:n]))
			*/
		}
	}
}
