package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const sample = `
3   4
4   3
2   5
1   3
3   9
3   3
`

func main() {
	var as []int
	var bs []int
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		a := mustInt(parts[0])
		as = append(as, a)
		b := mustInt(parts[1])
		bs = append(bs, b)
	}
	slices.Sort(as)
	slices.Sort(bs)

	fmt.Println(distance(as, bs))
	fmt.Println(similarity(as, bs))
}

func distance(as []int, bs []int) int {
	sum := 0
	for i := range as {
		sum = sum + abs(as[i]-bs[i])
	}
	return sum
}

func similarity(as []int, bs []int) int {
	ret := 0
	for _, a := range as {
		for _, b := range bs {
			if a == b {
				ret += a
			}
		}
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
