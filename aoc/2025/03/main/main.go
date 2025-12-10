package main

import (
	"fmt"
	"strings"
)

const sample = `
987654321111111
811111111111119
234234234234278
818181911112111
`

func main() {
	bat := parse(sample)
	fmt.Println(sumUp(bat, joltage1, true))
	fmt.Println(sumUp(bat, joltage2, true))
}

func sumUp(bat [][]int, joltage func([]int) int, debug bool) int {
	sum := 0
	for _, row := range bat {
		j := joltage(row)
		if debug {
			fmt.Println(j)
		}
		sum += j
	}
	return sum
}

func joltage1(row []int) int {
	hi := 0
	hiPos := 0
	for i, v := range row[:len(row)-1] {
		if v > hi {
			hiPos, hi = i, v
		}
	}

	lo := 0
	for _, v := range row[hiPos+1:] {
		if v > lo {
			lo = v
		}
	}
	return hi*10 + lo
}

func joltage2(row []int) int {
	return joltage(0, row, 12)
}

func joltage(acc int, row []int, sz int) int {
	if sz == 0 {
		return acc
	}
	sz--
	// find the highest number that leaves at least sz.
	val, pos := 0, 0
	for i, v := range row[:len(row)-sz] {
		if v > val {
			pos, val = i, v
		}
	}
	acc *= 10
	acc += val
	return joltage(acc, row[pos+1:], sz)
}

func parse(input string) [][]int {
	var ret [][]int
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var bat []int
		for _, c := range line {
			bat = append(bat, int(c-'0'))
		}
		ret = append(ret, bat)
	}
	return ret
}
