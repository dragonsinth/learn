package main

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const sample = `
1, 1
1, 6
8, 3
3, 4
5, 5
8, 9
`

const thresh = 32

var re = regexp.MustCompile(`^(\d+), (\d+)$`)

type pt struct {
	x, y int
}

func (p pt) dist(x int, y int) int {
	return abs(p.x-x) + abs(p.y-y)
}

func main() {
	var pts []pt

	lines := strings.Split(sample, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !re.MatchString(line) {
			panic(line)
		}
		m := re.FindStringSubmatch(line)
		pt := pt{
			x: mustInt(m[1]),
			y: mustInt(m[2]),
		}
		pts = append(pts, pt)
	}

	countsA := make([]int, len(pts))
	for y := 0; y <= 400; y++ {
		for x := 0; x <= 400; x++ {
			cl := closest(pts, x, y)
			if cl >= 0 {
				countsA[cl]++
			}
		}
	}
	fmt.Println(countsA)

	countsB := make([]int, len(pts))
	for y := -1; y <= 401; y++ {
		for x := -1; x <= 401; x++ {
			cl := closest(pts, x, y)
			if cl >= 0 {
				countsB[cl]++
			}
		}
	}
	fmt.Println(countsB)

	largest := 0
	for i, c := range countsA {
		if countsA[i] == countsB[i] {
			if c > largest {
				largest = c
			}
		}
	}
	fmt.Println(largest)

	// part 2
	count := 0
	for y := 0; y <= 400; y++ {
		for x := 0; x <= 400; x++ {
			sum := 0
			for _, pt := range pts {
				sum += pt.dist(x, y)
			}
			if sum < thresh {
				count++
			}
		}
	}
	fmt.Println(count)
}

func closest(pts []pt, x int, y int) int {
	type res struct {
		id   int
		dist int
	}
	var ret []res

	for id, pt := range pts {
		ret = append(ret, res{id: id, dist: pt.dist(x, y)})
	}

	slices.SortFunc(ret, func(a, b res) int {
		return a.dist - b.dist
	})

	if ret[0].dist < ret[1].dist {
		return ret[0].id
	}
	return -1
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
