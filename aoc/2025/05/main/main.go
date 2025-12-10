package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const sample = `
3-5
10-14
16-20
12-18

1
5
8
11
17
32
`

type rng struct {
	lo, hi int
}

func main() {
	rngs, ings := parse(sample)
	fmt.Println(countFresh(rngs, ings))
	fmt.Println(countAllFresh(rngs))
}

func countAllFresh(rngs []rng) int {
	slices.SortFunc(rngs, func(a, b rng) int {
		if a.lo != b.lo {
			return a.lo - b.lo
		}
		return a.hi - b.hi
	})

	var newRngs []rng
	for _, r := range rngs {
		if len(newRngs) > 0 && canMerge(newRngs[len(newRngs)-1], r) {
			newRngs[len(newRngs)-1] = merge(newRngs[len(newRngs)-1], r)
		} else {
			newRngs = append(newRngs, r)
		}
	}

	sum := 0
	for _, r := range newRngs {
		sum += r.hi - r.lo + 1
	}
	return sum
}

func canMerge(r1 rng, r2 rng) bool {
	// We know that r1 <= r2 on the left edge, which reduces the state space
	if r1.lo > r2.lo {
		panic("shouldn't happen")
	}

	return r2.lo <= r1.hi
}

func merge(r1 rng, r2 rng) rng {
	return rng{lo: min(r1.lo, r2.lo), hi: max(r1.hi, r2.hi)}
}

func countFresh(rngs []rng, ings []int) int {
	sum := 0
	for _, ing := range ings {
		if isFresh(ing, rngs) {
			sum++
		}
	}
	return sum
}

func isFresh(ing int, rngs []rng) bool {
	for _, rng := range rngs {
		if ing >= rng.lo && ing <= rng.hi {
			return true
		}
	}
	return false
}

func parse(input string) ([]rng, []int) {
	var rngs []rng
	var ings []int
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.ContainsRune(line, '-') {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				panic(parts)
			}
			rngs = append(rngs, rng{
				lo: mustInt(parts[0]),
				hi: mustInt(parts[1]),
			})
		} else {
			ings = append(ings, mustInt(line))
		}
	}
	return rngs, ings
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
