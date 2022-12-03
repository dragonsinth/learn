package main

import (
	"fmt"
	"strings"
)

var sample = `
vJrwpWtwJgWrhcsFMMfFFhFp
jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL
PmmdzqPrVvPwwTWBwg
wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn
ttgJtRGJQctTZtZT
CrZsJsPPZsGzwwsLwLmpwMDw
`

func main() {
	main1()

	sum := 0
	var lines []string

	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		lines = append(lines, line)
		if len(lines) < 3 {
			continue
		}

		// analyze
		reduced := make([]map[rune]bool, 3)
		for i, line := range lines {
			reduced[i] = map[rune]bool{}
			for _, c := range line {
				reduced[i][c] = true
			}
		}
		lines = nil

		counts := map[rune]int{}
		for _, r := range reduced {
			for k := range r {
				counts[k]++
			}
		}

		for k, v := range counts {
			if v >= 3 {
				fmt.Printf("%s -> %d\n", string(k), score(k))
				sum += score(k)
			}
		}

	}

	fmt.Println(sum)
}

func main1() {
	sum := 0

	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if len(line)%2 != 0 {
			panic("expect even line")
		}
		mid := len(line) / 2
		first := line[:mid]
		second := line[mid:]

		seen := map[rune]bool{}
		for _, c := range first {
			seen[c] = true
		}
		for _, c := range second {
			if seen[c] {
				fmt.Printf("%s %d\n", string(c), score(c))
				sum += score(c)
				break
			}
		}

	}

	fmt.Println(sum)
}

func score(c int32) int {
	switch {
	case c >= 'a' && c <= 'z':
		return int(c-'a') + 1
	case c >= 'A' && c <= 'Z':
		return int(c-'A') + 27
	default:
		panic(c)
	}
}
