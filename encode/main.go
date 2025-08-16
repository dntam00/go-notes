package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gogs/chardet"
)

func main() {
	encode()
}

func encode() {
	v := ""
	bytes := []byte(v)
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(bytes)
	if err != nil {
		fmt.Println("Error detecting encoding:", err)
		return
	}

	fmt.Printf("Detected encoding: %s\n", result.Charset)

	streamByte, _ := base64.URLEncoding.DecodeString(v)
	fmt.Println(string(streamByte))
}

func encode64() {

	sbData := ""
	dst := make([]byte, 3000)
	bytes := []byte(sbData)
	base64.URLEncoding.Encode(dst, bytes)

	fmt.Println(string(dst))
}
