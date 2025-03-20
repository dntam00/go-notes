package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"syscall"
	"time"
)

func setNonBlocking(fd int) error {
	// Get the current flags for the file descriptor
	flags, _, err := syscall.Syscall(syscall.SYS_FCNTL, uintptr(fd), syscall.F_GETFL, 0)
	if err != 0 {
		return fmt.Errorf("fcntl(F_GETFL) failed: %v", err)
	}

	// Set the O_NONBLOCK flag
	_, _, err = syscall.Syscall(syscall.SYS_FCNTL, uintptr(fd), syscall.F_SETFL, flags|syscall.O_NONBLOCK)
	if err != 0 {
		return fmt.Errorf("fcntl(F_SETFL) failed: %v", err)
	}

	return nil
}

func main() {
	// Create a TCP socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Println("Error creating socket:", err)
		return
	}
	defer syscall.Close(fd)

	// Set the socket to non-blocking mode
	if err := setNonBlocking(fd); err != nil {
		log.Println("Error setting non-blocking mode:", err)
		return
	}

	// Bind the socket to an address
	addr := syscall.SockaddrInet4{Port: 8085}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())

	if err := syscall.Bind(fd, &addr); err != nil {
		log.Println("Error binding socket:", err)
		return
	}

	// Listen for incoming connections
	if err := syscall.Listen(fd, 5); err != nil {
		log.Println("Error listening on socket:", err)
		return
	}

	log.Println("Server is listening on :8085")

	// Main event loop
	for {
		// Accept incoming connections (non-blocking)
		log.Println("start listen")
		nfd, _, err := syscall.Accept(fd)
		if err != nil {
			if errors.Is(err, syscall.EAGAIN) || errors.Is(err, syscall.EWOULDBLOCK) {
				// No incoming connections, continue polling
				time.Sleep(100 * time.Millisecond)
				continue
			}
			log.Println("Error accepting connection:", err)
			return
		}

		// Set the new socket to non-blocking mode
		if err := setNonBlocking(nfd); err != nil {
			log.Println("Error setting non-blocking mode for new connection:", err)
			syscall.Close(nfd)
			continue
		}

		log.Println("New connection accepted")

		// Handle the connection in a goroutine
		go handleConnection(nfd)
	}
}

func handleConnection(fd int) {
	defer syscall.Close(fd)

	buf := make([]byte, 1024)
	for {
		// Read data from the socket (non-blocking)
		n, err := syscall.Read(fd, buf)
		if err != nil {
			if errors.Is(err, syscall.EAGAIN) || errors.Is(err, syscall.EWOULDBLOCK) {
				// No data available, continue polling
				time.Sleep(100 * time.Millisecond)
				continue
			}
			log.Println("Error reading from socket:", err)
			return
		}

		if n == 0 {
			// Connection closed by client
			log.Println("Connection closed by client")
			return
		}

		// Process the received data
		message := string(buf[:n])
		log.Println("Received message:", message)

		// Echo the message back to the client
		if _, err := syscall.Write(fd, buf[:n]); err != nil {
			log.Println("Error writing to socket:", err)
			return
		}
	}
}
