package main

import (
	"github.com/songgao/water"
	"log"
)

func main() {
	ifce, err := water.New(water.Config{DeviceType: water.TUN})
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1500)
	for {
		n, err := ifce.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("packet: % x\n", buf[:n])
	}
}
