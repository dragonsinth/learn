package main

import (
	"fmt"
	"strconv"
	"strings"
)

const sample = `
+1
-2
+3
+1
`

func main() {
	seen := map[int]bool{0: true}
	sum := 0
	for {
		lines := strings.Split(sample, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			n := mustInt(line)
			sum += n
			if seen[sum] {
				fmt.Println(sum)
				return
			}
			seen[sum] = true
		}
	}
}

func main1() {
	sum := 0
	lines := strings.Split(sample, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n := mustInt(line)
		sum += n
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
