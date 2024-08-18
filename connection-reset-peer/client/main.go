package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	// Connect to the server on localhost, port 12345.
	conn, err := net.Dial("tcp4", "127.0.0.1:12345")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	fmt.Println("Connected to server")

	time.Sleep(7 * time.Second)

	_, err = conn.Write([]byte("Hello from client"))
	if err != nil {
		log.Println("write error: " + err.Error())
	}
	//
	//// sleep to make sure client receive "RST" packet
	time.Sleep(1 * time.Second)

	var bs = make([]byte, 1024)
	for {
		n, err := conn.Read(bs)
		if err != nil {
			log.Println("read messed up: " + err.Error())
			if err := conn.Close(); err != nil {
				log.Println("close connection error: " + err.Error())
			}
			break
		} else {
			fmt.Println("read", n, "bytes")
			fmt.Println("read: " + string(bs))
		}
		time.Sleep(1 * time.Second)
	}
}
