package main

import (
	"fmt"
	"net"
	"play-around/utils"
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

	err = conn.(*net.TCPConn).SetWriteBuffer(1)
	if err != nil {
		fmt.Println("error setting read buffer:", err)
	}

	go func() {
		for {
			buf := make([]byte, 1000)
			n, err := conn.Write(buf)
			time.Sleep(20 * time.Millisecond)
			if err != nil {
				fmt.Println("error writing:", err)
				continue
			}
			if n != 0 {
				fmt.Println("send n byte", n)
				continue
			}
			fmt.Println("write", n, "bytes")
		}
	}()

	utils.Wait()
	_ = conn.Close()
}
