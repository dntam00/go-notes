package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", ":9099")
	if err != nil {
		fmt.Println("Error setting up server:", err)
		return
	}
	fmt.Println("Server listening on port 9099")

	for {
		// Accept a connection.
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	var bs = make([]byte, 1024)
	var id string
	for {
		// set read timeout to make this read error and close connection
		_ = conn.SetDeadline(time.Now().Add(700 * time.Second))

		_ = conn.(*net.TCPConn).SetKeepAlive(false)

		_, err := conn.Read(bs)
		if err != nil {
			_ = conn.Close()
			if id == "1" {
				log.Println("closed at " + time.Now().String())
			}
			break
		} else {
			id = "1"
			fmt.Println("read: " + string(bs) + " at " + time.Now().String())
		}
	}
}
