package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/grid"
	"github.com/dragonsinth/learn/aoc/termbox"
	"golang.org/x/exp/slices"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var sample = `
R 6 (#70c710)
D 5 (#0dc571)
L 2 (#5713f0)
D 2 (#d2c081)
R 2 (#59c680)
D 2 (#411b91)
L 5 (#8ceee2)
U 2 (#caa173)
L 1 (#1b58a2)
U 2 (#caa171)
R 2 (#7807d2)
U 3 (#a77fa3)
L 2 (#015232)
U 2 (#7a21e3)
`

func main() {
	run(sample, false, true)
	run(sample, true, false)
}

func run(input string, part2 bool, debug bool) {
	p := parse(input, part2)
	if debug {
		fmt.Println("start")
		termbox.RenderPlain(p.Render(), os.Stdout)
	}
	fmt.Println(p.Count())
}

var re = regexp.MustCompile(`^([LRUD]) (\d+) \(#([0-9a-f]+)\)$`)

func parse(input string, part2 bool) puz {
	var p puz
	lines := strings.Split(input, "\n")
	x, y := 0, 0
	for _, line := range lines {
		if line == "" {
			continue
		}
		if !re.MatchString(line) {
			panic(line)
		}
		matches := re.FindStringSubmatch(line)

		var d dir
		var dist int
		if !part2 {
			d = dirMap[matches[1][0]]
			dist = mustInt(matches[2])
			_, err := strconv.ParseInt(matches[3], 16, 32)
			if err != nil {
				panic(matches[3])
			}
		} else {
			hx := matches[3]
			d = dir(hx[5] - '0')
			parseDist, err := strconv.ParseInt(matches[3][:5], 16, 32)
			if err != nil {
				panic(matches[3])
			}
			dist = int(parseDist)
		}

		switch d {
		case R:
			x += dist
		case D:
			p.cols = append(p.cols, col{
				x:    x,
				ymin: y,
				ymax: y + dist,
			})
			y += dist
		case L:
			x -= dist
		case U:
			p.cols = append(p.cols, col{
				x:    x,
				ymin: y - dist,
				ymax: y,
			})
			y -= dist
		}

		p.xmin = min(p.xmin, x)
		p.xmax = max(p.xmax, x)
		p.ymin = min(p.ymin, y)
		p.ymax = max(p.ymax, y)
	}
	slices.SortFunc(p.cols, func(a, b col) bool {
		if a.x != b.x {
			return a.x < b.x
		}
		return a.ymin < b.ymin
	})
	return p
}

type col struct {
	x          int
	ymin, ymax int
}

type puz struct {
	cols       []col
	xmin, xmax int
	ymin, ymax int
}

func (p puz) Render() [][]byte {
	buf := grid.Alloc2d(p.xmax-p.xmin+1, p.ymax-p.ymin+1, byte(' '))
	p.trace(func(x1, x2 int, y int, wat byte) {
		for x := x1; x <= x2; x++ {
			buf[y-p.ymin][x-p.xmin] = wat
		}
	})
	return buf
}

func (p puz) Count() int {
	sum := 0
	p.trace(func(x1, x2 int, y int, wat byte) {
		if wat != '.' {
			sum += (x2 - x1) + 1
		}
	})
	return sum
}

func (p puz) trace(fill func(x1, x2 int, y int, wat byte)) {
	for y := p.ymin; y <= p.ymax; y++ {
		// !in, !rowEdge = fully outside
		// in, !rowEdge = fully inside
		// in, rowEdge = a top edge we're inside
		// !in, rowEdge = a bottom edge we're outside
		in := false
		rowEdge := false

		commit := func(x1, x2 int) {
			switch {
			case !in && !rowEdge:
				fill(x1, x2, y, '.')
			case in && !rowEdge:
				fill(x1, x2, y, '#')
			case in && rowEdge:
				fill(x1, x2, y, '=')
			case !in && rowEdge:
				fill(x1, x2, y, '.')
			default:
				panic("here")
			}
		}

		x := p.xmin
		for _, c := range p.cols {
			if y < c.ymin {
				// ignore this column
			} else if y == c.ymin {
				// handle top edges
				if in {
					if rowEdge {
						// was a top edge we're terminating
						fill(x, c.x-1, y, '=')
						x = c.x
						fill(x, x, y, '7')
					} else {
						// start a bottom edge
						fill(x, c.x-1, y, '#')
						x = c.x
						fill(x, x, y, 'F')
					}
				} else {
					if rowEdge {
						// was a bottom edge we're terminating
						fill(x, c.x-1, y, '_')
						x = c.x
						fill(x, x, y, '7')
					} else {
						// start a top edge
						fill(x, c.x-1, y, '.')
						x = c.x
						fill(x, x, y, 'F')
					}
				}
				in = !in
				rowEdge = !rowEdge
				x++
			} else if y < c.ymax {
				// could be starting or ending a run or an edge
				commit(x, c.x-1)
				x = c.x
				fill(x, x, y, '|')
				x++
				in = !in
			} else if y == c.ymax {
				// handle bottom edges
				if rowEdge {
					if in {
						fill(x, c.x-1, y, '=') // was a top edge we're terminating
					} else {
						fill(x, c.x-1, y, '_') // was a bottom edge we're terminating
					}
					x = c.x
					fill(x, x, y, 'J')
					rowEdge = false
					// in is unaffected by hovering below an edge
				} else {
					// starting a bottom edge
					if in {
						fill(x, c.x-1, y, '#') // was a top edge we're terminating
					} else {
						fill(x, c.x-1, y, '.') // was a bottom edge we're terminating
					}
					x = c.x
					fill(x, x, y, 'L')
					rowEdge = true
					// in is unaffected by hovering below an edge
				}
				x++
			} else {
				// ignore this column
			}
		}
		// better not be inside
		if in || rowEdge {
			panic("unexpected")
		}
		fill(x, p.xmax, y, '.')
	}
}

type dir byte

const (
	R = dir(iota)
	D
	L
	U
)

var dirMap = map[byte]dir{
	'R': R,
	'D': D,
	'L': L,
	'U': U,
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
