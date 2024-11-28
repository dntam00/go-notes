package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/sys/unix"
	"math/big"
	"net"
	"os"
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

// https://stackoverflow.com/questions/22116873/set-socket-option-is-why-so-important-for-a-socket-ip-hdrincl-in-icmp-request
// https://darkcoding.net/software/raw-sockets-in-go-link-layer/

type IcmpResult struct {
	dst         string
	result      string
	seq         int
	requestTime time.Time
}

// go run main.go 5 20.1.0.12 veth1
func main() {

	argsWithoutProg := os.Args[1:]

	numsOfPackets, _ := strconv.Atoi(argsWithoutProg[0])
	destination := argsWithoutProg[1]

	bindInterface := ""
	if len(argsWithoutProg) > 2 {
		bindInterface = argsWithoutProg[2]
	}

	timeout := 4 * time.Second
	var err error
	sendingFd, _ := unix.Socket(unix.AF_INET, unix.SOCK_RAW, unix.IPPROTO_RAW)

	if bindInterface != "" {
		err = unix.SetsockoptString(sendingFd, unix.SOL_SOCKET, unix.SO_BINDTODEVICE, bindInterface)
		if err != nil {
			panic("failed to bind fd to interface, error: " + err.Error())
		}
	}

	responseCh := make(chan IcmpResult, numsOfPackets)
	done := make(chan bool)
	packets := make(map[int]IcmpResult)

	go func() {
		for {
			select {
			case result, closed := <-responseCh:
				if !closed {
					done <- true
					break
				}
				origin := packets[result.seq]
				fmt.Printf("ping to: %v, response: %v, seq: %d, latency:%f\n", origin.dst, result.result, result.seq, time.Since(origin.requestTime).Seconds())
			}
		}
	}()

	icmpFd, _ := unix.Socket(unix.AF_INET, unix.SOCK_RAW, unix.IPPROTO_ICMP)
	icmpFdFile := os.NewFile(uintptr(icmpFd), fmt.Sprintf("fd %d", icmpFd))

	defer func() {
		err := icmpFdFile.Close()
		if err != nil {
			fmt.Println("failed to close icmp socket, error: " + err.Error())
		}
		_ = unix.Close(icmpFd)
	}()

	go func(f *os.File) {
		received := 0
		for {
			buf := make([]byte, 256)
			tv := syscall.Timeval{
				Sec:  int64(timeout.Seconds()),
				Usec: int64(timeout.Microseconds() % 1e6),
			}
			now := time.Now()
			if err := syscall.SetsockoptTimeval(icmpFd, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv); err != nil {
				panic("failed to set read deadline, error: " + err.Error())
			}
			packetLen, readErr := f.Read(buf)
			if readErr != nil {
				fmt.Println("read timeout, take: ", time.Since(now).Milliseconds())
				continue
			}

			ipHeader := ipv4.Header{}
			readErr = ipHeader.Parse(buf[:packetLen])

			if readErr != nil {
				fmt.Println("failed to parse ip ipHeader, error: " + readErr.Error())
				continue
			}

			icmpPkt := buf[ipHeader.Len:ipHeader.TotalLen]

			icmpMessage, readErr := icmp.ParseMessage(1, icmpPkt)
			if readErr != nil {
				fmt.Println("failed to parse icmp icmpMessage, error: " + readErr.Error())
				continue
			}

			//fmt.Printf("receive icmp packet: %v, %v, %d\n", icmpMessage.Type, icmpMessage.Code, packetLen)

			// assume v4
			// https://datatracker.ietf.org/doc/html/rfc792
			//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			//|     Type      |     Code      |          Checksum             |
			//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			//|                             unused                            |
			//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			//|      Internet Header + 64 bits of Original Data Datagram      |
			//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

			icmpConcrete := icmpMessage.Type.(ipv4.ICMPType)

			switch icmpConcrete {
			case ipv4.ICMPTypeDestinationUnreachable:
				internalPayload, readErr := icmpMessage.Body.Marshal(1)
				if readErr != nil {
					fmt.Println("failed to parse icmp icmp internal payload, error: " + readErr.Error())
					continue
				}
				// skip 4 unused bytes
				internalPayload = internalPayload[4:]

				originalIPHeader := ipv4.Header{}
				readErr = originalIPHeader.Parse(internalPayload)

				if readErr != nil {
					fmt.Println("failed to parse original ip ipHeader, error: " + readErr.Error())
					continue
				}
				if originalIPHeader.Dst.String() != destination {
					fmt.Println("destination mismatch, skip destination unreachable packet")
					continue
				}

				received++

				// skip 6 bytes of original data to get seq number
				seqStr := internalPayload[originalIPHeader.Len+6 : originalIPHeader.Len+8]
				seqNumber := new(big.Int)
				seqNumber.SetString(fmt.Sprintf("%X", seqStr), 16)

				responseCh <- IcmpResult{result: "destination unreachable", seq: int(seqNumber.Uint64())}
				break
			case ipv4.ICMPTypeEchoReply:
				if ipHeader.Src.String() != destination {
					fmt.Println("destination mismatch, skip echo reply packet")
					continue
				}
				received++
				internalPayload, readErr := icmpMessage.Body.Marshal(1)
				if readErr != nil {
					fmt.Println("failed to parse icmp icmp internal payload, error: " + readErr.Error())
					continue
				}
				seqStr := internalPayload[2:4]
				seqNumber := new(big.Int)
				seqNumber.SetString(fmt.Sprintf("%X", seqStr), 16)
				responseCh <- IcmpResult{result: "ping reply", seq: int(seqNumber.Int64())}
				break
			case ipv4.ICMPTypeEcho:
				continue
			}
			if received == numsOfPackets {
				close(responseCh)
				break
			}
		}
	}(icmpFdFile)

	for seq := 1; seq <= numsOfPackets; seq++ {
		packets[seq] = IcmpResult{dst: destination, seq: seq, requestTime: time.Now()}
		p := buildPacket(destination, seq)

		addr := unix.SockaddrInet4{
			Port: 0,
			Addr: [4]byte(net.ParseIP(destination).To4()),
		}
		err = unix.Sendto(sendingFd, p, 0, &addr)
		if err != nil {
			panic("failed to send packet, error: " + err.Error())
		}
	}

	_ = <-done
}

func buildPacket(dest string, seq int) []byte {
	ipHeader := ipv4.Header{
		Version: ipv4.Version,
		Len:     20,
		TOS:     0,
		//TotalLen: 20 + 10,
		//ID:       0,
		Flags:    ipv4.DontFragment,
		FragOff:  0,
		TTL:      32,
		Protocol: 1,
		//Checksum: 0,
		//Src: net.ParseIP("20.1.0.10"),
		Dst: net.ParseIP(dest),
		//Options: nil,
	}

	echo := &icmp.Echo{
		ID:   1,
		Seq:  seq,
		Data: []byte{0xC0, 0xDE},
	}

	icmpMsg := icmp.Message{
		Type:     ipv4.ICMPTypeEcho,
		Code:     0,
		Checksum: 0,
		Body:     echo,
	}

	t, err := icmpMsg.Marshal(nil)

	out, err := ipHeader.Marshal()
	if err != nil {
		panic("failed to marshal ip header, error: " + err.Error())
	}
	return append(out, t...)
}
