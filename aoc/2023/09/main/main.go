package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `
0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45
`

func main() {
	sum := 0
	psum := 0
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		val := computeNext(parseNums(line))
		fmt.Println(val)
		sum += val

		pval := computePrev(parseNums(line))
		fmt.Println(pval)
		psum += pval
	}
	fmt.Println(sum)
	fmt.Println(psum)
}

func computeNext(nums []int) int {
	//fmt.Println(nums)

	// base case
	if allZeroes(nums) {
		return 0
	}

	diffs := make([]int, len(nums)-1)
	for i := range diffs {
		diffs[i] = nums[i+1] - nums[i]
	}

	next := computeNext(diffs)
	return nums[len(nums)-1] + next
}

func computePrev(nums []int) int {
	// fmt.Println(nums)

	// base case
	if allZeroes(nums) {
		return 0
	}

	diffs := make([]int, len(nums)-1)
	for i := range diffs {
		diffs[i] = nums[i+1] - nums[i]
	}

	next := computePrev(diffs)
	return nums[0] - next
}

func allZeroes(nums []int) bool {
	for _, v := range nums {
		if v != 0 {
			return false
		}
	}
	return true
}

func parseNums(s string) []int {
	var ret []int
	for _, p := range strings.Fields(s) {
		ret = append(ret, mustInt(p))
	}
	return ret
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
