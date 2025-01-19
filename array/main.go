package main

import "fmt"

func main() {

	//ints := make([]int, 2)
	//ints = append(ints, 0)
	//fmt.Println(len(ints))

	//access()
	loop()
}

type data struct {
	value int
}

func access() {
	arr := make([]data, 1)
	arr[0] = data{value: 1}
	before := &arr[0]
	fmt.Println(before.value)
	arr[0] = data{value: 2}
	fmt.Println(before.value)
}

type Char struct {
	Char *string
}

func loop() {
	var chars []Char
	values := []string{"a", "b", "c"}
	for _, v := range values {
		fmt.Println("v: ", &v)
		chars = append(chars, Char{Char: &v})
	}
	for _, v := range chars {
		fmt.Println(*v.Char)
	}
}
