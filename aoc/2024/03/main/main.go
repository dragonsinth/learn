package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	sample  = `xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))`
	sample2 = `xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))`
)

var (
	re = regexp.MustCompile(`mul\(([0-9]{1,3}),([0-9]{1,3})\)`)
)

func main() {
	fmt.Println(part1(sample))
	fmt.Println(part2(sample2))
}

func part1(in string) int {
	sum := 0
	matches := re.FindAllString(in, -1)
	for _, m := range matches {
		subs := re.FindStringSubmatch(m)
		a := mustInt(subs[1])
		b := mustInt(subs[2])
		sum += a * b
	}
	return sum
}

var (
	outerRe = regexp.MustCompile(`(mul\([0-9]{1,3},[0-9]{1,3}\))|(do\(\))|(don't\(\))`)
)

func part2(in string) int {
	sum := 0
	enabled := true
	matches := outerRe.FindAllString(in, -1)
	for _, m := range matches {
		if m == `do()` {
			enabled = true
		} else if m == `don't()` {
			enabled = false
		} else if strings.HasPrefix(m, `mul(`) {
			subs := re.FindStringSubmatch(m)
			a := mustInt(subs[1])
			b := mustInt(subs[2])
			if enabled {
				sum += a * b
			}
		}
	}
	return sum
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
