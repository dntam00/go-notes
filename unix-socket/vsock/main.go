package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	vsock "github.com/Code-Hex/darwin-vsock"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

// see `man vsock`
func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "err: %+v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	l, err := vsock.Listen("vsock", &vsock.Addr{CID: unix.VMADDR_CID_HOST, Port: 4321})
	if err != nil {
		return errors.WithStack(err)
	}
	defer l.Close()

	go func() {
		conn, err := l.Accept()
		if err != nil {
			log.Println("accept error", err)
			return
		}
		go func() {
			defer conn.Close()
			_, err := io.Copy(conn, conn)
			if err != nil {
				log.Println("copy err", err)
			}
		}()
	}()

	log.Println("listener addr:", l.Addr())

	conn, err := vsock.DialContext(ctx, "vsock", nil, &vsock.Addr{CID: unix.VMADDR_CID_HYPERVISOR, Port: 4321})
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	msg := "hello, world"
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return errors.WithStack(err)
	}

	result := make([]byte, len(msg))
	n, err := conn.Read(result)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Printf("successful: %q\n", string(result[:n]))

	return nil
}
