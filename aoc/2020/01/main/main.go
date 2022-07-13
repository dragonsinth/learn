package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `
1721
979
366
299
675
1456
`

func main() {
	var vals []int
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		vals = append(vals, mustInt(line))
	}

	for i, a := range vals {
		for j, b := range vals {
			if i < j && a+b == 2020 {
				fmt.Println(a, b, a*b)
			}
		}
	}

	for i, a := range vals {
		for j, b := range vals {
			for k, c := range vals {
				if i < j && j < k && a+b+c == 2020 {
					fmt.Println(a, b, c, a*b*c)
				}
			}
		}
	}

}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
