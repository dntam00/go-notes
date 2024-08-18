package main

import (
	"crypto/tls"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var addr = flag.String("addr", "", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{
		Scheme: "wss",
		Host:   *addr, Path: "/w/",
		RawQuery: "token=",
	}
	log.Printf("connecting to %s", u.String())

	dialer := websocket.DefaultDialer
	c, _, err := dialer.Dial(u.String(), nil)

	err = c.UnderlyingConn().(*tls.Conn).NetConn().(*net.TCPConn).SetKeepAlivePeriod(2 * time.Second)

	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	_ = c.WriteJSON(map[string]interface{}{
		"type": "ping",
		"time": 123,
	})
	time.Sleep(40 * time.Second)
	c.WriteJSON(map[string]interface{}{
		"type": "ping",
		"time": 123,
	})
	time.Sleep(80 * time.Second)
}
