package main

import (
	"fmt"
	"strconv"
	"strings"
)

const sample = `
L68
L30
R48
L5
R60
L55
L1
L99
R14
L82
`

func main() {
	deltas := parse(sample)
	part1(deltas, true)
	part2(deltas, true)
}

func parse(input string) []int {
	var deltas []int
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		switch line[0] {
		case 'L':
			deltas = append(deltas, -mustInt(line[1:]))
		case 'R':
			deltas = append(deltas, mustInt(line[1:]))
		default:
			panic(line)
		}
	}
	return deltas
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

func part1(deltas []int, debug bool) {
	pos := 50
	count := 0
	for _, delta := range deltas {
		next := mod(pos+delta, 100)
		if debug {
			fmt.Printf("The dial is rotated %d to point at %d\n", delta, next)
		}
		pos = next
		if pos == 0 {
			count++
		}
	}
	fmt.Println(count)
}

func part2(deltas []int, debug bool) {
	pos := 50
	count := 0
	for _, delta := range deltas {
		zeroes := 0
		next := pos + delta
		if next < 0 {
			zeroes += -next / 100 // e.g. -250 would generate 2 complete spins
			if pos != 0 {
				// if we weren't already on 0, we crossed it to go negative
				zeroes++
			}
			next = mod(next, 100)
		} else if next >= 100 {
			zeroes += next / 100
			next = mod(next, 100)
		} else if next == 0 {
			zeroes++
		}
		if debug {
			fmt.Printf("The dial is rotated %d to point at %d, generating %d zeroes\n", delta, next, zeroes)
		}
		pos = next
		count += zeroes
	}
	fmt.Println(count)
}

func mod(a int, m int) int {
	return (a%m + m) % m
}
