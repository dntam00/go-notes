package main

import (
	"fmt"
	"log"
	"net"
	"syscall"
)

const (
	MaxEvents      = 10
	ReadBufferSize = 1024
)

func main() {
	// Create epoll instance
	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(epfd)

	// Create a TCP listener
	listener, err := net.Listen("tcp", ":8085")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// Get the file descriptor for the listener
	listenerFd := int(listener.(*net.TCPListener).File().Fd())

	// Add the listener file descriptor to epoll
	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(listenerFd),
	}
	if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, listenerFd, &event); err != nil {
		log.Fatal(err)
	}

	// Create buffer for events
	events := make([]syscall.EpollEvent, MaxEvents)

	// Event loop
	for {
		n, err := syscall.EpollWait(epfd, events, -1)
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < n; i++ {
			if events[i].Fd == int32(listenerFd) {
				// Accept new connections
				conn, err := listener.Accept()
				if err != nil {
					log.Println("Error accepting connection:", err)
					continue
				}

				// Get the file descriptor for the connection
				connFd := int(conn.(*net.TCPConn).File().Fd())

				// Add the connection file descriptor to epoll
				event := syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(connFd),
				}
				if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, connFd, &event); err != nil {
					log.Println("Error adding connection to epoll:", err)
					conn.Close()
					continue
				}

				log.Println("New connection accepted")
			} else if events[i].Events&syscall.EPOLLIN != 0 {
				// Read data from the connection
				buffer := make([]byte, ReadBufferSize)
				n, err := syscall.Read(int(events[i].Fd), buffer)
				if err != nil {
					log.Println("Error reading from connection:", err)
					syscall.Close(int(events[i].Fd))
					continue
				}
				if n == 0 {
					// Connection closed by client
					log.Println("Connection closed by client")
					syscall.Close(int(events[i].Fd))
					continue
				}

				// Process the received data
				message := string(buffer[:n])
				fmt.Printf("Read %d bytes: %s\n", n, message)

				// Echo the message back to the client
				if _, err := syscall.Write(int(events[i].Fd), buffer[:n]); err != nil {
					log.Println("Error writing to connection:", err)
					syscall.Close(int(events[i].Fd))
				}
			}
		}
	}
}
