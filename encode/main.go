package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gogs/chardet"
)

func main() {
	v := ""
	bytes := []byte(v)
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(bytes)
	if err != nil {
		fmt.Println("Error detecting encoding:", err)
		return
	}

	fmt.Printf("Detected encoding: %s\n", result.Charset)

	streamByte, _ := base64.StdEncoding.DecodeString(v)
	fmt.Println(string(streamByte))
}
