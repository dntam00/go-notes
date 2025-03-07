package main

import (
	"fmt"
	"math"
)

func main() {
	//fmt.Println(minSpeedOnTime([]int{1, 1, 100000}, 2.01))
	fmt.Println(leftMost([]int{1}, 2))
}

func leftMost(arr []int, target int) int {
	left := 0
	right := len(arr)
	for !(left >= right) {
		mid := (left + right) / 2
		if arr[mid] == target {
			right = mid
		} else if arr[mid] > target {
			right = mid
		} else if arr[mid] < target {
			left = mid + 1
		}
	}
	if left == len(arr) || arr[left] != target {
		return -1
	}
	return left
}

func minSpeedOnTime(dist []int, hour float64) int {
	var total float64
	total = 0.0

	lenDist := len(dist) - 1

	if float64(lenDist) >= hour {
		return -1
	}

	left := 1
	high := math.MaxInt32

	//lastResult := -1

	// [1,2,3,4,5]
	// get both left and high

	for left < high {
		middle := (left + high) / 2
		fmt.Println(middle)
		for j := 0; j < lenDist; j++ {
			total += math.Ceil(float64(dist[j]) / float64(middle))
		}
		total += float64(dist[lenDist]) / float64(middle)
		if total <= hour {
			// keep middle in the result
			high = middle
		} else {
			left = middle + 1
		}
		total = 0.00
	}

	//for j := 0; j < lenDist; j++ {
	//	total += math.Ceil(float64(dist[j]) / float64(high))
	//}
	//total += float64(dist[lenDist]) / float64(high)
	//if total > hour {
	//	return -1
	//}

	return high
}
