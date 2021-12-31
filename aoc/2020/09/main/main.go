package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var sample = `
35
20
15
25
47
40
62
55
65
95
102
117
150
182
127
219
299
277
309
576
`

var (
	input = sample
)

const preamble = 25

func main() {
	var nums []int

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		nums = append(nums, mustInt(line))
	}

	spooky := 0
	for i, n := range nums {
		if i < preamble {
			continue
		}
		if !validate(nums[i-preamble:i], n) {
			fmt.Println(i, n)
			spooky = n
			break
		}
	}

	for i := range nums {
		if vals := isSpooky(spooky, nums[i:]); vals != nil {
			fmt.Println(vals, vals[0] + vals[len(vals) - 1])
			break
		}
	}
}

func isSpooky(spooky int, ints []int) []int {
	sum := 0
	for i, n := range ints {
		sum += n
		switch {
		case sum < spooky:
			continue
		case sum > spooky:
			return nil
		case sum == spooky:
			ret := append([]int{}, ints[:i+1]...)
			sort.Ints(ret)
			return ret
		default:
			panic("wat")
		}
	}
	return nil
}

func validate(ints []int, n int) bool {
	for i := range ints {
		for j := range ints {
			if i != j && ints[i]+ints[j] == n {
				return true
			}
		}
	}

	return false
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
