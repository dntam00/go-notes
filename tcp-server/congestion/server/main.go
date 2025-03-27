package main

import (
	"fmt"
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
	// don't consume data
	_ = conn.SetDeadline(time.Now().Add(700 * time.Second))
	_ = conn.(*net.TCPConn).SetKeepAlive(false)
	//err := conn.(*net.TCPConn).SetReadBuffer(9000)
	//if err != nil {
	//	fmt.Println("error setting read buffer:", err)
	//	return
	//}
	select {}
}
