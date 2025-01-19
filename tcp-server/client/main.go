package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp4", ":9099")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	_ = conn.(*net.TCPConn).SetKeepAlive(false)

	_, _ = conn.Write([]byte("Hello from lb"))

	go func() {
		time.Sleep(400 * time.Second)
		if _, err = conn.Write([]byte("Hello from lb after 400 seconds")); err != nil {
			fmt.Println("Error writing after 400 seconds:", err)
		}
	}()

	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			break
		} else {
			fmt.Println("Received:", string(buf))
		}
		if err = conn.Close(); err != nil {
			fmt.Println("Error closing connection:", err)
		}
	}
}
