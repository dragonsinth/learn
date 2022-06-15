package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var sample = `
0,9 -> 5,9
8,0 -> 0,8
9,4 -> 3,4
2,2 -> 2,1
7,0 -> 7,4
6,4 -> 2,0
0,9 -> 2,9
3,4 -> 1,4
0,0 -> 8,8
5,5 -> 8,2
`

const width = 10

var input = sample

type puzzle struct {
	count [width * width]int
}

func (p *puzzle) Print() {
	for y := 0; y < width; y++ {
		fmt.Println(p.count[y*width : y*width+width])
	}
}

func (p *puzzle) Counts() int {
	sum := 0
	for _, v := range p.count {
		if v > 1 {
			sum++
		}
	}
	return sum
}

func (p *puzzle) Inc(x, y int) {
	p.count[y*width+x]++
}

var parseLine = regexp.MustCompile(`^(\d+),(\d+)\s*->\s*(\d+),(\d+)$`)

func main() {
	var puz puzzle
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !parseLine.MatchString(line) {
			panic(line)
		}

		subs := parseLine.FindStringSubmatch(line)
		if len(subs) != 5 {
			panic(subs)
		}
		fmt.Println(subs)

		x1, y1, x2, y2 := mustInt(subs[1]), mustInt(subs[2]), mustInt(subs[3]), mustInt(subs[4])
		if x1 == x2 {
			step := 1
			if y2 < y1 {
				step = -1
			}
			for y := y1; ; y += step {
				puz.Inc(x1, y)
				if y == y2 {
					break
				}
			}
		} else if y1 == y2 {
			step := 1
			if x2 < x1 {
				step = -1
			}
			for x := x1; ; x += step {
				puz.Inc(x, y1)
				if x == x2 {
					break
				}
			}
		} else {
			// Must be a diagonal.
			ystep := 1
			if y2 < y1 {
				ystep = -1
			}
			xstep := 1
			if x2 < x1 {
				xstep = -1
			}
			for x, y := x1, y1; ; x, y = x+xstep, y+ystep {
				puz.Inc(x, y)
				if x == x2 && y == y2 {
					break
				} else if x == x2 || y == y2 {
					panic("not diagonal")
				}
			}
		}

	}

	puz.Print()
	fmt.Println(puz.Counts())
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
