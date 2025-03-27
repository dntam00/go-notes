package main

import (
	"log"
	"math"
)

func main() {
	u := uint32(0)
	u = math.MaxUint32
	log.Println(u + 1)
}
