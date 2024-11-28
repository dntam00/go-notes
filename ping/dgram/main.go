package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	target := "20.1.0.11" // Replace with your target IP
	ifaceName := "veth1"  // Replace with your desired interface name

	// Resolve the target address
	destAddr := &net.IPAddr{IP: net.ParseIP(target)}
	if destAddr == nil {
		log.Fatalf("Invalid target address: %s", target)
	}

	// Create a socket with SOCK_DGRAM and IPPROTO_ICMP
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_ICMP)
	if err != nil {
		log.Fatalf("Failed to create raw socket: %v", err)
	}
	defer syscall.Close(fd)

	// Bind the socket to a specific network interface
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Fatalf("Failed to get interface %s: %v", ifaceName, err)
	}
	err = syscall.SetsockoptString(fd, syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, iface.Name)
	if err != nil {
		log.Fatalf("Failed to bind socket to device %s: %v", ifaceName, err)
	}

	// Set up the ICMP Echo Request
	bytes := make([]byte, 3000)
	c := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: bytes,
		},
	}

	// Marshal the ICMP message
	messageBytes, err := c.Marshal(nil)
	if err != nil {
		log.Fatalf("Failed to marshal ICMP message: %v", err)
	}

	// Send the packet
	start := time.Now()
	destSockAddr := &syscall.SockaddrInet4{}
	copy(destSockAddr.Addr[:], destAddr.IP.To4())

	err = syscall.Sendto(fd, messageBytes, 0, destSockAddr)
	if err != nil {
		log.Fatalf("Failed to send ICMP packet: %v", err)
	}

	for {
		buf := make([]byte, 1500)
		syscall.SetsockoptTimeval(fd, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &syscall.Timeval{Sec: 3, Usec: 0})

		n, from, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			log.Fatalf("Failed to receive ICMP response: %v", err)
		}

		// Parse the ICMP response
		response, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), buf[:n])
		if err != nil {
			log.Fatalf("Failed to parse ICMP response: %v", err)
		}

		switch response.Type {
		case ipv4.ICMPTypeEchoReply:
			elapsed := time.Since(start)
			fmt.Printf("Ping to %s through %s: seq=%d time=%v, from=%+v\n",
				target, ifaceName, 1, elapsed, from)
		case ipv4.ICMPTypeDestinationUnreachable:
			fmt.Printf("Destination unreachable, with code: %d\n", response.Code)

		default:
			fmt.Printf("Unexpected response: %+v\n", response)
		}
	}
}
