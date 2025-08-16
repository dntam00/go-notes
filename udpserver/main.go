package main

import (
	"fmt"
	"net"
)

func main() {
	addr := &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 5000}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()
	fmt.Println("UDP server listening on", addr)

	buf := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Read error:", err)
			continue
		}
		fmt.Printf("Received from %v: %s\n", clientAddr, string(buf[:n]))

		// Send response to client
		_, err = conn.WriteToUDP([]byte("Hello from server"), clientAddr)
		if err != nil {
			fmt.Println("Write error:", err)
		}
	}
}
