package main

import "log"

type Api struct {
}

func (a *Api) Func1() {
	log.Println("Func 1")
}

type ApiDecorator struct {
	*Api
}

func (a *ApiDecorator) Func1Decorated() {
	log.Println("Decorator Func 1")
}

func main() {
	api := &Api{}
	api.Func1()

	decoratedApi := &ApiDecorator{Api: api}
	decoratedApi.Func1() // Calls the original Func1 from Api
}
