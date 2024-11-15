package main

import "fmt"

//- export GitTags=$(git describe --tags)
//- export RevVersion=$(git rev-parse --short HEAD)
//- export BuildTime=$(date -u "+%Y%m%d%I%M")
//- go build -ldflags "-X main.AppVersion=$GitTags.$RevVersion.$BuildTime" -o bin/linker cmd/linker/main.go

var Tag string

func main() {
	fmt.Println(Tag)
}
