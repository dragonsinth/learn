package main

import (
	"fmt"
	"strings"
)

func part1(input string, w, h int, turns int) {
	wm, hm := w/2, h/2
	qs := [4]int{}
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if !re.MatchString(line) {
			panic(line)
		}
		matches := re.FindStringSubmatch(line)
		if len(matches) != 5 {
			panic(matches)
		}

		x, y, dx, dy := mustInt(matches[1]), mustInt(matches[2]), mustInt(matches[3]), mustInt(matches[4])

		fx := mod(x+dx*turns, w)
		fy := mod(y+dy*turns, h)
		if fx < wm && fy < hm {
			qs[0]++
		} else if fx > wm && fy < hm {
			qs[1]++
		} else if fx < wm && fy > hm {
			qs[2]++
		} else if fx > wm && fy > hm {
			qs[3]++
		} else if fx != wm && fy != hm {
			panic("here")
		}
	}
	fmt.Println(qs)
	fmt.Println(qs[0] * qs[1] * qs[2] * qs[3])
}

func mod(n, m int) int {
	r := n % m
	if r < 0 {
		r += m
	}
	if r < 0 {
		panic(r)
	}
	return r
}
