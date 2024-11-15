package main

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"syscall"
)

func main() {
	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 36890})
	_, err := conn.WriteTo([]byte("hello"), &net.UDPAddr{IP: net.ParseIP(""), Port: 12345})
	if err != nil {
		fmt.Println("send 1: ", err)
	}
	_, err = conn.WriteTo([]byte("hello"), &net.UDPAddr{IP: net.ParseIP(""), Port: 12345})
	if err != nil {
		fmt.Println("send 2: ", err)
	}
	fmt.Println("finish")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
}
