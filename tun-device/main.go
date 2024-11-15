package main

import (
	"fmt"
	"golang.org/x/net/ipv4"
	"log"
	"os"
	"syscall"
	"unsafe"
)

// Constants for TUN device setup
const (
	TUN_DEVICE = "/dev/net/tun"
	IFF_TUN    = 0x0001
	IFF_NO_PI  = 0x1000
)

// OpenTUN opens a TUN device and returns a file descriptor
func OpenTUN(ifaceName string) (*os.File, error) {
	// Open the TUN device
	file, err := os.OpenFile(TUN_DEVICE, os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open TUN device: %v", err)
	}

	// Set up the TUN device with the specified name
	var ifr [18]byte
	copy(ifr[:], ifaceName)
	*(*uint16)(unsafe.Pointer(&ifr[16])) = IFF_TUN | IFF_NO_PI

	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		file.Fd(),
		uintptr(syscall.TUNSETIFF),
		uintptr(unsafe.Pointer(&ifr[0])),
	)
	if errno != 0 {
		return nil, fmt.Errorf("failed to set up TUN interface: %v", errno)
	}

	return file, nil
}

// ReadPacket reads an IP packet from the TUN interface
func ReadPacket(tun *os.File) {
	buf := make([]byte, 1500) // Buffer for reading packets

	for {
		// Read a packet from the TUN interface
		n, err := tun.Read(buf)
		if err != nil {
			log.Fatalf("failed to read from TUN interface: %v", err)
		}

		// Parse the IP packet
		packet := buf[:n]
		header, err := ipv4.ParseHeader(packet)
		if err != nil {
			log.Printf("failed to parse packet: %v", err)
			continue
		}

		// Print packet details
		fmt.Printf("Received packet from %s to %s, Protocol: %d\n",
			header.Src, header.Dst, header.Protocol)
	}
}

func main() {
	// Open and configure the TUN interface
	tun, err := OpenTUN("tun0")
	if err != nil {
		log.Fatalf("Error opening TUN interface: %v", err)
	}
	defer func() {
		_ = tun.Close()
	}()

	// Read packets from the TUN interface
	fmt.Println("Reading packets from TUN interface...")
	ReadPacket(tun)
}
