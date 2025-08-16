package main

import (
	"errors"
	"fmt"
	"time"
)

func main() {
	//deferClosure()
	//testError()
	testError2()
}

func deferClosure() {
	var i int
	defer func() {
		fmt.Print(i)
	}()
	for ; i < 5; i++ {

	}
}

type Err struct {
	val int
}

func testError2() {

	//var err *Err
	//fmt.Println(err.val)
	//err = &Err{val: 1}

	var err2 Err

	err2 = Err{
		val: 1,
	}

	defer func(err2 *Err) {
		fmt.Println(err2.val)
	}(&err2)

	err2 = Err{val: 2}
}

func testError() {
	var err error
	go func(er error) {
		time.Sleep(5 * time.Second)
		fmt.Println(er)
	}(err)

	err = errors.New("hmm")
	time.Sleep(6 * time.Second)
}
