package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"
)

func main() {
	// Create a new file descriptor for stdin
	fd, err := syscall.Dup(syscall.Stdin)
	if err != nil {
		fmt.Printf("Failed to duplicate stdin: %v\n", err)
		return
	}
	defer syscall.Close(fd)

	// Set non-blocking mode
	if err := syscall.SetNonblock(fd, true); err != nil {
		fmt.Printf("Failed to set non-blocking mode: %v\n", err)
		return
	}

	fmt.Println("Type something (program will check for input every second):")

	buffer := make([]byte, 1024)
	for {
		// Try to read from stdin (non-blocking)
		n, err := syscall.Read(fd, buffer)
		if err != nil {
			if errors.Is(err, syscall.EAGAIN) || errors.Is(err, syscall.EWOULDBLOCK) {
				// No data available right now
				fmt.Println("receive signal from os:", err, errors.Is(err, syscall.EAGAIN), errors.Is(err, syscall.EWOULDBLOCK))
				time.Sleep(5 * time.Second)
				continue
			}
			fmt.Printf("Error reading from stdin: %v\n", err)
			return
		}

		if n > 0 {
			// Print the received input
			input := string(buffer[:n])
			fmt.Printf("Received input: %s", input)
			if input == "quit\n" {
				fmt.Println("Exiting...")
				os.Exit(0)
			}
		}
	}
}
