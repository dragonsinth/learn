package main

import (
	"fmt"
	"strings"
)

var sample = `
mjqjpqmgbljsphdztnvjfqwrcgsmlb
bvwbjplbgvbhsrlpgdmjqwftvncz
nppdvjthqldpwncqszvftbrmjlhg
nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg
zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw
`

func main() {
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		for i := 4; i < len(line); i++ {
			if unique(line[i-4 : i]) {
				fmt.Println(i)
				break
			}
		}

		for i := 14; i < len(line); i++ {
			if unique(line[i-14 : i]) {
				fmt.Println(i)
				break
			}
		}
	}

	fmt.Println()
}

func unique(s string) bool {
	seen := [256]bool{}
	for _, c := range s {
		if seen[c] {
			return false
		}
		seen[c] = true
	}
	return true
}
