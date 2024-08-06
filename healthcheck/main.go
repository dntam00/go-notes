package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	CheckTcpPort("127.0.0.1", "7996", 60)
}

func CheckTcpPort(host string, port string, timeoutSecond int) {
	timeout := time.Second * time.Duration(timeoutSecond)
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		fmt.Println("connect failed: ", err)
	}
	if conn != nil {
		defer func(conn net.Conn) {
			err := conn.Close()
			if err != nil {
				fmt.Printf("Error close tcp connection: %v", err)
			}
		}(conn)
	}
}
