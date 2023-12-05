package main

import (
	"fmt"
	"strings"
)

var sample1 = `
1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet
`

var sample2 = `
two1nine
eightwothree
abcone2threexyz
xtwone3four
4nineeightseven2
zoneight234
7pqrstsixteen
`

var nums = []string{
	"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
}

func main() {
	part1(sample1)
	part2(sample2)
}

func part1(input string) {
	sum := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var digits []int
		for _, c := range line {
			if c >= '0' && c <= '9' {
				digits = append(digits, int(c-'0'))
			}
		}
		f, l := digits[0], digits[len(digits)-1]
		val := f*10 + l
		sum += val
	}
	fmt.Println(sum)
}

func part2(input string) {
	sum := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var digits []int
		for i, c := range line {
			if c >= '0' && c <= '9' {
				digits = append(digits, int(c-'0'))
			}

			// check line for exact suffix
			prev := line[:i+1]
			for val, test := range nums {
				if strings.HasSuffix(prev, test) {
					digits = append(digits, val)
				}
			}
		}
		f, l := digits[0], digits[len(digits)-1]
		val := f*10 + l
		sum += val
	}
	fmt.Println(sum)
}
