package main

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

func main() {
	fmt.Println(maximumSum([]int{229, 398, 269, 317, 420, 464, 491, 218, 439, 153, 482, 169, 411, 93, 147, 50, 347, 210, 251, 366, 401}))
}

func clearDigits(s string) string {
	var clearF func(input string) string
	clearF = func(input string) string {
		for i, c := range input {
			if unicode.IsDigit(c) {
				j := i - 1
				for j >= 0 {
					if !unicode.IsDigit(rune(input[j])) {
						break
					}
					j--
				}
				if j == -1 {
					return clearF(input[i+1:])
				}
				return clearF(input[0:j] + input[i+1:])
			}
		}
		return input
	}
	return clearF(s)
}

func clearStrings(s string, part string) string {
	var clearF func(input, p string) string
	clearF = func(input, p string) string {
		for i, _ := range input {
			if strings.HasPrefix(input[i:], p) {
				return clearF(input[0:i]+input[i+len(p):], part)
			}
		}
		return input
	}
	return clearF(s, part)
}

type pair struct {
	val       int
	sumDigits int
}

func maximumSum(nums []int) int {
	pairs := make([]pair, len(nums))
	for i, v := range nums {
		pairs[i] = pair{
			val:       v,
			sumDigits: calculateSum(v),
		}
	}

	m := make(map[int][]int)
	for _, p := range pairs {
		ns, ok := m[p.sumDigits]
		if !ok {
			m[p.sumDigits] = []int{p.val}
		}
		ns = append(ns, p.val)
		m[p.sumDigits] = ns
	}

	maxSum := -1

	for _, v := range m {
		if len(v) >= 2 {
			sort.Ints(v[:])
			currentSum := v[len(v)-1] + v[len(v)-2]
			if currentSum > maxSum {
				maxSum = currentSum
			}
		}
	}

	return maxSum

	//sort.Slice(pairs, func(i, j int) bool {
	//	return pairs[i].sumDigits < pairs[j].sumDigits
	//})
	//for _, v := range pairs {
	//	fmt.Println(v.sumDigits)
	//}
	//maxSum := -1
	//for i := 0; i < len(pairs)-1; i++ {
	//	if pairs[i].sumDigits == pairs[i+1].sumDigits {
	//		currentSum := pairs[i].val + pairs[i+1].val
	//		if currentSum > maxSum {
	//			maxSum = currentSum
	//		}
	//	}
	//}
}

func calculateSum(num int) int {
	sum := 0
	for num > 0 {
		sum += num % 10
		num /= 10
	}
	return sum
}
