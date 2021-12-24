package main

import (
	"fmt"
	"sort"
)

func main() {
	main1()
	fmt.Println()

	if false {
		survived := map[int64]bool{}
		for v := int64(189); v <= 405; v++ {
			// Which first six steps can survive the last 2?
			if anySolution([]int{}, 0, 0, 0, 0, v, 6) {
				survived[v] = true
			}
		}

		// Which can survive the second 6?
		old := toSlice(survived)
		survived = map[int64]bool{}
		for _, v := range old {
			if anySolution([]int{}, 0, 0, v, 6, v, 6) {
				survived[v] = true
			}
		}

		// Which survive the last 2?
		old = toSlice(survived)
		survived = map[int64]bool{}
		for _, v := range old {
			if anySolution([]int{}, 0, 0, v, 12, 0, 2) {
				survived[v] = true
				fmt.Println(v)
			}
		}
	}

	code := []int64{6, 2, 9, 1, 1, 9, 4, 1, 7, 1, 6, 1, 1, 1}

	var x, y, z int64
	for i, c := range code {
		if i == 4 || i == 6 || i == 12 {
			fmt.Println()
		}
		_, x, y, z = steps[i](int64(c), x, y, z)
		fmt.Printf("%2d: %d -> (%d,%d,%d)\n", i, c, x, y, z)
	}
}

func anySolution(choices []int, x, y, z int64, step int, expect int64, depth int) bool {
	if depth == 0 {
		if z == expect {
			fmt.Println(expect, ":", choices)
			return true
		}
		return false
	}
	for i := 1; i <= 9; i++ {
		_, x, y, z := steps[step](int64(i), x, y, z)
		if anySolution(append(choices, i), x, y, z, step+1, expect, depth-1) {
			return true
		}
	}
	return false
}

func toSlice(in map[int64]bool) []int64 {
	ret := []int64{}
	for k := range in {
		ret = append(ret, k)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i] < ret[j]
	})
	return ret
}

func main1() {
	code := []int64{9, 9, 9, 1, 1, 9, 9, 3, 9, 4, 9, 6, 8, 4}

	var x, y, z int64
	for i, c := range code {
		if i == 4 || i == 6 || i == 12 {
			fmt.Println()
		}
		_, x, y, z = steps[i](int64(c), x, y, z)
		fmt.Printf("%2d: %d -> (%d,%d,%d)\n", i, c, x, y, z)
	}
}
