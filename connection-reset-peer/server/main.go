package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println("Error setting up server:", err)
		return
	}
	fmt.Println("Server listening on port 12345")

	for {
		// Accept a connection.
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Client connected")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	var bs = make([]byte, 1024)
	for {
		// set read timeout to make this read error and close connection
		_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
		n, err := conn.Read(bs)
		if err != nil {
			log.Println("read messed up: " + err.Error())
			_ = conn.Close()
			break
		} else {
			fmt.Println("read", n, "bytes")
			fmt.Println("read: " + string(bs))
		}
		time.Sleep(time.Second)
	}
	fmt.Println("Connection closed by server")
}
