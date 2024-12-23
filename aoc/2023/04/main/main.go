package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var sample = `
Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19
Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1
Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83
Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36
Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11
`

func main() {
	sum, cardSum := 0, 0
	acc := make([]int, 1024)

	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			panic(parts)
		}
		id := mustInt(strings.Fields(parts[0])[1])
		acc[id]++

		split := strings.Split(parts[1], " | ")
		winning := parseNums(split[0])
		ours := parseNums(split[1])

		matches := 0
		for _, n := range ours {
			if slices.Index(winning, n) >= 0 {
				matches++
			}
		}

		// score
		score := 0
		if matches > 0 {
			score = 1 << (matches - 1)
		}

		// accumulate
		nCopies := acc[id]
		for i := 0; i < matches; i++ {
			acc[id+i+1] += nCopies
		}

		fmt.Println("id", id, "=", score, acc[id])
		sum += score
		cardSum += acc[id]
	}
	fmt.Println(sum, cardSum)
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
