package main

import "math"

func main() {
	b := New([]int{5, 1, 1})
	for i := 0; i < 10; i++ {
		println(b.Next())
	}
}

type balancer struct {
	weighted           []int
	currentState       []int
	totalWeight        int
	currentDistributed int
}

func New(weighted []int) *balancer {
	totalWeighted := 0
	for _, c := range weighted {
		totalWeighted += c
	}
	return &balancer{
		weighted:           weighted,
		currentState:       make([]int, len(weighted)),
		totalWeight:        totalWeighted,
		currentDistributed: 0,
	}
}

func (b *balancer) Next() int {
	for i, s := range b.currentState {
		b.currentState[i] = s + b.weighted[i]
	}
	selected := -1
	selectedState := math.MinInt
	for i, s := range b.currentState {
		if s > selectedState {
			selected = i
			selectedState = s
		}
	}
	b.currentState[selected] -= b.totalWeight
	b.currentDistributed++
	if b.currentDistributed == b.totalWeight {
		b.currentDistributed = 0
		b.currentState = make([]int, len(b.weighted))
	}
	return selected
}
