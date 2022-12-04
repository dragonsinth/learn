package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `
2-4,6-8
2-3,4-5
5-7,7-9
2-8,3-7
6-6,4-6
2-6,4-8
`

func main() {
	contained := 0
	overlap := 0

	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			panic("expect 2 parts")
		}

		a, b := parseRange(parts[0]), parseRange(parts[1])
		if a.contains(b) || b.contains(a) {
			contained++
		}

		if a.overlap(b) {
			overlap++
		}
	}

	fmt.Println(contained, overlap)
}

type rng struct {
	lo, hi int
}

func (a rng) contains(b rng) bool {
	return a.lo <= b.lo && a.hi >= b.hi
}

func (a rng) overlap(b rng) bool {
	return a.lo <= b.hi && a.hi >= b.lo
}

func parseRange(s string) rng {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		panic("expect 2 parts")
	}
	lo := mustInt(parts[0])
	hi := mustInt(parts[1])
	if lo > hi {
		panic(s)
	}
	return rng{lo: lo, hi: hi}
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
