package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	host string
	port string
}

func main() {
	server := Server{
		host: "127.0.0.1",
		port: "7995",
	}
	server.Run()
}

// Run ...
func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = listener.Close()
	}()

	log.Printf("Server started at port %v", server.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(conn net.Conn) {
			defer func() {
				_ = conn.Close()
			}()
			fmt.Println("Connection accepted")
		}(conn)
	}
}
