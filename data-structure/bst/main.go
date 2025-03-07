package main

import "fmt"

type node struct {
	value int
	left  *node
	right *node
}

type bst struct {
	root *node
}

func (b *bst) init() {
	b.root = nil
}

func (b *bst) add(value int) {
	b.root = doAdd(value, b.root)
}

func doAdd(value int, current *node) *node {
	if current == nil {
		return &node{
			value: value,
		}
	}
	if current.value > value {
		current.left = doAdd(value, current.left)
	}
	if current.value < value {
		current.right = doAdd(value, current.right)
	}
	return current
}

func main() {
	b := &bst{}
	b.add(5)
	b.add(4)
	b.add(7)
	b.add(6)
	b.add(2)
	b.add(8)
	b.add(9)
	fmt.Println(b.root)
}
