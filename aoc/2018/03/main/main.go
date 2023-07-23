package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const sample = `
#1 @ 1,3: 4x4
#2 @ 3,1: 4x4
#3 @ 5,5: 2x2
`

var re = regexp.MustCompile(`^#(\d+) @ (\d+),(\d+): (\d+)x(\d+)$`)

func main() {
	counts := make([][]int, 1000)
	for i := range counts {
		counts[i] = make([]int, 1000)
	}

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
		_, l, t, w, h := mustInt(m[1]), mustInt(m[2]), mustInt(m[3]), mustInt(m[4]), mustInt(m[5])

		for y, ym := t, t+h; y < ym; y++ {
			for x, xm := l, l+w; x < xm; x++ {
				counts[y][x]++
			}
		}
	}

	var overlap int
	for _, row := range counts {
		for _, c := range row {
			if c > 1 {
				overlap++
			}
		}
	}
	fmt.Println(overlap)

	// part 2
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !re.MatchString(line) {
			panic(line)
		}
		m := re.FindStringSubmatch(line)
		id, l, t, w, h := mustInt(m[1]), mustInt(m[2]), mustInt(m[3]), mustInt(m[4]), mustInt(m[5])

		maybe := true
		for y, ym := t, t+h; y < ym; y++ {
			for x, xm := l, l+w; x < xm; x++ {
				if counts[y][x] > 1 {
					maybe = false
				}
			}
		}
		if maybe {
			fmt.Println(id)
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
