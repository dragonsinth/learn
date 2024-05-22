package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `
12
14
1969
100756
`

func main() {
	sum := 0
	lines := strings.Split(sample, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n := mustInt(line)
		for {
			n = n/3 - 2
			if n > 0 {
				sum += n
			} else {
				break
			}
		}
	}
	fmt.Println(sum)
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
