package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var sample = `
1000
2000
3000

4000

5000
6000

7000
8000
9000

10000
`

func main() {
	var elves []int
	cur := 0

	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			if cur != 0 {
				elves = append(elves, cur)
				cur = 0
			}
			continue
		}
		cur += mustInt(line)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(elves)))
	fmt.Printf("The top elf has %d\n", elves[0])
	fmt.Printf("The top 3 elves have %d\n", elves[0]+elves[1]+elves[2])

}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
