package main

import (
	"errors"
	"net"
	"strconv"
	"time"
)

var errFailedToCastAddr = errors.New("failed to cast net.Addr to *net.UDPAddr or *net.TCPAddr")

func main() {
	conn, _ := net.ListenPacket("udp4", "0.0.0.0:0")
	addr := conn.LocalAddr()
	ip, port, _ := AddrIPPort(addr)
	println(ConvertString(ip, port))
	time.Sleep(50 * time.Second)
}

func AddrIPPort(a net.Addr) (net.IP, int, error) {
	aUDP, ok := a.(*net.UDPAddr)
	if ok {
		return aUDP.IP, aUDP.Port, nil
	}

	aTCP, ok := a.(*net.TCPAddr)
	if ok {
		return aTCP.IP, aTCP.Port, nil
	}

	return nil, 0, errFailedToCastAddr
}

func ConvertString(ip net.IP, port int) string {
	return net.JoinHostPort(ip.String(), strconv.Itoa(port))
}
