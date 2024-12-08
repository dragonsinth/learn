package main

import (
	"fmt"
	"strconv"
	"strings"
)

const sample = `
190: 10 19
3267: 81 40 27
83: 17 5
156: 15 6
7290: 6 8 6 15
161011: 16 10 13
192: 17 8 14
21037: 9 7 18 13
292: 11 6 16 20
`

func main() {
	run(false, sample)
	run(true, sample)
}

func run(allowConcat bool, in string) {
	sum := 0
	for _, line := range strings.Split(in, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		val := mustInt(parts[0])
		ints := parseNums(parts[1])
		if canEqual(allowConcat, val, ints) {
			sum += val
		}
	}
	fmt.Println(sum)
}

func canEqual(allowConcat bool, val int, ints []int) bool {
	if len(ints) == 1 {
		return ints[0] == val
	}
	last := ints[len(ints)-1]

	// try division (reverse mult)
	if (val % last) == 0 {
		if canEqual(allowConcat, val/last, ints[:len(ints)-1]) {
			return true
		}
	}

	// try subtraction (reverse add)
	if val > last {
		if canEqual(allowConcat, val-last, ints[:len(ints)-1]) {
			return true
		}
	}

	// try substring (reverse concat)
	if allowConcat {
		// quick test
		if val%10 == last%10 {
			valStr := strconv.Itoa(val)
			lastStr := strconv.Itoa(last)
			if strings.HasSuffix(valStr, lastStr) {
				vStr := strings.TrimSuffix(valStr, lastStr)
				v := 0
				if len(vStr) > 0 {
					v = mustInt(vStr)
				}
				if canEqual(allowConcat, v, ints[:len(ints)-1]) {
					return true
				}
			}
		}
	}

	return false
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
