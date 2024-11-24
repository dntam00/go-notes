package main

import (
	"context"
	"fmt"
	"golang.org/x/net/ipv4"
	"golang.org/x/sys/unix"
	"math/big"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

//+0------7-------15---------------31
//|  Type | Code  |    Checksum    |
//+--------------------------------+
//| Iden  |    Sequence Number     |
//+--------------------------------+
//|             Data               |
//+--------------------------------+

func main() {

	numsOfPackets := 5
	destination := "20.1.0.11"
	bindInterface := "veth1"
	timeout := 10 * time.Second

	var err error
	fd, _ := unix.Socket(unix.AF_INET, unix.SOCK_RAW, unix.IPPROTO_RAW)

	if bindInterface != "" {
		err = unix.SetsockoptString(fd, unix.SOL_SOCKET, unix.SO_BINDTODEVICE, bindInterface)
		if err != nil {
			panic("failed to bind fd to interface, error: " + err.Error())
		}
	}

	go func() {
		fd, _ := unix.Socket(unix.AF_INET, unix.SOCK_RAW, unix.IPPROTO_ICMP)
		//unix.SetsockoptString(fd, unix.SOL_SOCKET, unix.SO_BINDTODEVICE, "veth1")
		f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))

		for {
			buf := make([]byte, 256)
			tv := syscall.Timeval{
				Sec:  int64(timeout.Seconds()),
				Usec: int64(timeout.Microseconds() % 1e6),
			}
			now := time.Now()
			if err := syscall.SetsockoptTimeval(fd, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv); err != nil {
				panic("failed to set read deadline, error: " + err.Error())
			}
			_, err := f.Read(buf)
			if err != nil {
				fmt.Println("read timeout, take: ", time.Since(now).Milliseconds())
				return
			}
			icmpTypeStr := buf[20:21]
			icmpType, _ := strconv.Atoi(fmt.Sprintf("% X", icmpTypeStr))

			if icmpType == 3 {
				fmt.Printf("receive destination unreachable, currentSeq: % X", buf)
			}
			if icmpType == 0 {
				seqStr := buf[26:28]
				seqNumber := new(big.Int)
				seqNumber.SetString(fmt.Sprintf("%X", seqStr), 16)
				sendTimeByte := buf[28:43]
				sendTime := time.Time{}
				_ = sendTime.UnmarshalBinary(sendTimeByte)
				fmt.Printf("receive icmp reply packet: seq=%d, latency= %dμs\n", seqNumber, time.Since(sendTime).Microseconds())
			}
		}
	}()

	for seq := 1; seq <= numsOfPackets; seq++ {
		p := buildPacket(destination, seq)

		addr := unix.SockaddrInet4{
			Port: 0,
			Addr: [4]byte(net.ParseIP(destination).To4()),
		}
		err = unix.Sendto(fd, p, 0, &addr)
		if err != nil {
			panic("failed to send packet, error: " + err.Error())
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
}

func buildPacket(dest string, seq int) []byte {
	ipHeader := ipv4.Header{
		Version:  4,
		Len:      20,
		TOS:      0,
		TotalLen: 20 + 10,
		//ID:       0,
		Flags:    0,
		FragOff:  0,
		TTL:      32,
		Protocol: 1,
		//Checksum: 0,
		//Src: net.ParseIP("20.1.0.10"),
		Dst: net.ParseIP(dest),
		//Options: nil,
	}

	icmpPacket := []byte{
		8, // type: echo request
		0, // code: not used by echo request
		0, // checksum (16 bit), we fill in below
		0,
		0, // identifier (16 bit). zero allowed.
		1,
		//byte(seq), // sequence number (16 bit). zero allowed.
		//0,
		////0xC0, // Optional data. ping puts time packet sent here
		//...now,
		//0xDE,
	}

	seqBitStr := strconv.FormatInt(int64(seq), 2)
	for len(seqBitStr) <= 16 {
		seqBitStr = "0" + seqBitStr
	}
	seqFirstByteStr := seqBitStr[:len(seqBitStr)-8]
	seqSecondByteStr := seqBitStr[len(seqBitStr)-8:]
	seqFirstByte, _ := strconv.ParseInt(seqFirstByteStr, 2, 64)
	seqSecondByte, _ := strconv.ParseInt(seqSecondByteStr, 2, 64)
	icmpPacket = append(icmpPacket, byte(seqFirstByte))
	icmpPacket = append(icmpPacket, byte(seqSecondByte))

	// timestamp len: 15
	now, err := time.Now().MarshalBinary()
	icmpPacket = append(icmpPacket, now...)

	icmpPacket = append(icmpPacket, 0xDE)
	cs := csum(icmpPacket)
	icmpPacket[2] = byte(cs)
	icmpPacket[3] = byte(cs >> 8)

	out, err := ipHeader.Marshal()
	if err != nil {
		panic("failed to marshal ipHeader, error: " + err.Error())
	}
	return append(out, icmpPacket...)
}

func csum(b []byte) uint16 {
	var s uint32
	for i := 0; i < len(b); i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	// add back the carry
	s = s>>16 + s&0xffff
	s = s + s>>16
	return uint16(^s)
}
