package main

import (
	"fmt"
	"strings"
)

var sample = `
BFFFBBFRRR
FFFBBBFRRR
BBFFBBFRLL
`

var (
	input = sample
)

func main() {
	max := 0
	seen := [1024]bool{}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		val := 0
		for _, b := range []byte(line) {
			val <<= 1
			switch b {
			case 'F', 'L':
				// nothing
			case 'B', 'R':
				val += 1
			default:
				panic(b)
			}
		}

		seen[val] = true
		if val > max {
			max = val
		}
	}
	fmt.Println(max)
	for i, v := range seen {
		if !v {
			fmt.Println(i)
		}
	}
}
