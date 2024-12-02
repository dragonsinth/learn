package main

import (
	"fmt"
	"strconv"
	"strings"
)

const sample = `
7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9
`

func main() {
	safe := 0
	lenientSafe := 0
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		vals := parseNums(line)
		if isSafe(vals) {
			safe++
			lenientSafe++
		} else if isLenientSafe(vals) {
			lenientSafe++
		}
	}
	fmt.Println(safe, lenientSafe)
}

func isLenientSafe(vals []int) bool {
	for i := range vals {
		if isSafe(without(vals, i)) {
			return true
		}
	}
	return false
}

func without(vals []int, i int) []int {
	ret := append([]int{}, vals[:i]...)
	ret = append(ret, vals[i+1:]...)
	return ret
}

func isSafe(vals []int) bool {
	increasing := vals[0] < vals[1]
	for i := range vals {
		if i == 0 {
			continue
		}
		prev := vals[i-1]
		curr := vals[i]
		if increasing != (prev < curr) {
			return false
		}
		diff := abs(prev - curr)
		if diff < 1 || diff > 3 {
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
