package main

import (
	"fmt"
	"reflect"
)

type AppError struct {
	code      int
	errorType string
	cause     error
}

type CError struct {
}

func (e *CError) Error() string {
	return "error A"
}

type Dog struct {
	name string
}

func (d *Dog) SetName(name string) {
	d.name = name
}

func (d *Dog) Name() string {
	return d.name
}

type Pet interface {
	Name() string
}

func main() {
	//var x interface{}
	//var y *int = nil
	//x = y
	//
	//if x != nil {
	//	fmt.Println("x != nil") // actual
	//} else {
	//	fmt.Println("x == nil") // expect
	//}
	//
	//fmt.Println(x)
	//test()

	//typeOf()

	//testInterface()
	//checkAssignError()

	testC()
}

func testC() {
	if WrapError() == nil {
		println("nil")
	} else {
		println("not nil")
	}
}

func test() {
	var err error
	err = getErrorInterface()
	if err == nil {
		println("nil")
	} else {
		println("not nil")
	}
}

func testInterface() {
	dog := Dog{name: "hello"}

	var pet Pet = &dog
	dog.SetName("hi")
	fmt.Println(reflect.ValueOf(pet))
	fmt.Println(dog.name)
	fmt.Println(reflect.TypeOf(pet))
	fmt.Println(getType(pet))
	//fmt.Println(reflect.TypeOf((Pet)(pet)).Elem().String())
}

func getType(x interface{}) string {
	return reflect.TypeOf(x).String()
}

func typeOf() {
	var err error
	fmt.Println(reflect.TypeOf(err))
	err = getErr()
	fmt.Println(reflect.TypeOf(err))
}

func WrapError() error {
	return getErr()
}

func getErr() *CError {
	return nil
}

func getErrorInterface() error {
	return nil
}

func checkAssignError() {
	var err error
	err = getErr()
	if err == nil {
		fmt.Println("nil")
	} else {
		fmt.Println("not nil")
	}
	fmt.Println(reflect.TypeOf(err))
	fmt.Println(reflect.ValueOf(err).IsNil())
}
