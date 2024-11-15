package main

import "fmt"

func main() {

	str := "ana"
	fmt.Println(len(str)) // 3
	str = "世界"
	fmt.Println(len(str)) // 6 not 2

	var i []int
	if len(i) == 0 {
		println("empty")
	} else {
		println("not empty")
	}
}
