package main

import (
	"fmt"
	"regexp"
	"strconv"
)

const sample = `
p=0,4 v=3,-3
p=6,3 v=-1,-3
p=10,3 v=-1,2
p=2,0 v=2,-1
p=0,0 v=1,3
p=3,0 v=-2,-2
p=7,6 v=-1,-3
p=3,0 v=-1,-2
p=9,3 v=2,3
p=7,3 v=-1,2
p=2,4 v=2,-3
p=9,5 v=-3,-3

`

var (
	re = regexp.MustCompile(`^p=(\d+),(\d+) v=(-?\d+),(-?\d+)$`)
)

// redacted!!
const (
	wRepeat    = 99
	hRepeat    = 99
	p2solution = 99
)

func main() {
	part1(sample, 11, 7, 100)
	// part1(data, 101, 103, 100)

	part2(sample, 11, 7, 77, true)
	// part2(data, 101, 103, p2solution, false)
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
