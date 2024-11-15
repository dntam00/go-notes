package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a Unix domain socket and listen for incoming connections.
	socket, err := net.Listen("unix", "/tmp/echo.sock")
	if err != nil {
		log.Fatal(err)
	}

	target, err := net.Listen("unix", "/tmp/echo_target.sock")

	go func() {
		for {
			// Accept an incoming connection.
			conn, err := target.Accept()
			if err != nil {
				log.Fatal(err)
			}

			// Handle the connection in a separate goroutine.
			go func(conn net.Conn) {
				defer conn.Close()
				// Create a buffer for incoming data.

				for {
					buf := make([]byte, 4096)

					// Read data from the connection.
					n, err := conn.Read(buf)
					if err != nil {
						log.Fatal(err)
					}

					log.Println("target receive: ", string(buf[:n]))
				}
			}(conn)
		}
	}()

	targetConn, err := net.Dial("unix", "/tmp/echo_target.sock")

	// Cleanup the sockfile.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove("/tmp/echo.sock")
		os.Exit(1)
	}()

	for {
		// Accept an incoming connection.
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Handle the connection in a separate goroutine.
		go func(sourceConn net.Conn) {
			defer sourceConn.Close()
			// Create a buffer for incoming data.

			_, err2 := io.Copy(targetConn, sourceConn)
			if err2 != nil {
				fmt.Println("error copy: ", err2)
			}
		}(conn)
	}
}
