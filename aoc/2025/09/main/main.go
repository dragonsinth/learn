package main

import (
	"fmt"
	"strconv"
	"strings"
)

const sample = `
7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3
`

type pos struct {
	x, y int
}

func (p pos) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

func main() {
	p := parse(sample)
	p.part1()
	p.part2()
}

func parse(input string) *puz {
	p := &puz{}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			panic(len(parts))
		}
		p.pts = append(p.pts, pos{mustInt(parts[0]), mustInt(parts[1])})
	}
	return p
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

func area(a, b pos) int {
	return dist(a.x, b.x) * dist(a.y, b.y)
}

func dist(a, b int) int {
	if a < b {
		return b - a + 1
	}
	return a - b + 1
}

func sort(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}
