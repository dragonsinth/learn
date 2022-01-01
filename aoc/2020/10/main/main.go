package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var sample = `
16
10
15
5
1
11
7
19
6
12
4
`

var sample2 = `
28
33
18
42
31
14
46
20
48
47
24
23
49
45
19
38
39
11
1
32
25
35
8
17
7
9
4
2
34
10
3
`

var (
	input = sample
)

func main() {
	nums := []int{0}

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		nums = append(nums, mustInt(line))
	}

	sort.Ints(nums)
	nums = append(nums, nums[len(nums)-1]+3)
	fmt.Println(nums)
	diffs := [4]int{}
	for i := range nums {
		if i == 0 {
			continue
		}
		diff := nums[i] - nums[i-1]
		diffs[diff]++
	}

	fmt.Println(diffs, diffs[1]*diffs[3])

	fmt.Println()

	// Break the chain into segments bounded by gaps of 3, multiply through.
	prod := 1
	last := 0
	for i := range nums {
		if i-last > 5 && nums[i]-nums[i-1] == 3 {
			prod *= countPerms(nums[last:i])
			last = i
		}
	}
	// Do the final segment.
	prod *= countPerms(nums[last:])
	fmt.Println(prod)
}

func countPerms(nums []int) int {
	return doCountPerms(nil, nums[0], nums[1:])
}

func doCountPerms(path []int, num int, nums []int) int {
	path = append(path, num)
	if len(nums) == 0 {
		fmt.Println(path)
		return 1
	}
	sum := 0
	for i := range nums {
		if num+3 >= nums[i] {
			sum += doCountPerms(path, nums[i], nums[i+1:])
		} else {
			break
		}
	}
	return sum
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
