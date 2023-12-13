package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"strings"
)

var sample = `
#.##..##.
..#.##.#.
##......#
##......#
..#.##.#.
..##..##.
#.#.##.#.

#...##..#
#....#..#
..##..###
#####.##.
#####.##.
..##..###
#....#..#
`

func main() {
	run(sample, 0, true)
	run(sample, 1, true)
}

func run(sample string, expectDiffs int, debug bool) {
	ps := parse(sample)
	sum := 0
	for i, p := range ps {
		if debug {
			termbox.RenderPlain(p.Render(), os.Stdout)
		}

		// compute folds
		for x := 1; x < p.w; x++ {
			if p.DiffsBeforeCol(x) == expectDiffs {
				fmt.Println(i, "cols", x)
				sum += x
			}
		}
		for y := 1; y < p.h; y++ {
			if p.DiffsFoldBeforeRow(y) == expectDiffs {
				fmt.Println(i, "rows", y)
				sum += 100 * y
			}
		}
	}
	fmt.Println(sum)
}

func parse(input string) []puz {
	var ps []puz
	var p *puz
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if line == "" {
			if p != nil {
				ps = append(ps, *p)
				p = nil
			}
			continue
		}
		if p == nil {
			p = &puz{
				pts: map[pos]wat{},
				w:   len(line),
				h:   0,
			}

		} else if p.w != len(line) {
			panic(line)
		}

		for x, c := range []byte(line) {
			w := wat(c)
			if w == '#' {
				pt := pos{x, p.h}
				p.pts[pt] = w
			} else if w != '.' {
				panic(c)
			}
		}
		p.h++
	}

	if p != nil {
		ps = append(ps, *p)
		p = nil
	}
	return ps
}

type puz struct {
	pts  map[pos]wat
	w, h int
}

func (p puz) at(pt pos) wat {
	v, ok := p.pts[pt]
	if !ok {
		return '.'
	}
	return v
}

func (p puz) Render() [][]byte {
	var buf [][]byte
	for y := 0; y < p.h; y++ {
		var line []byte
		for x := 0; x < p.w; x++ {
			w := p.at(pos{x, y})
			line = append(line, byte(w))
		}
		buf = append(buf, line)
	}
	return buf
}

func (p puz) DiffsBeforeCol(x int) int {
	// x represents the col after the fold, so x=2, w big: folds 1 onto 2
	// x=7, w=10: 6..7 folds onto 8..9
	diffs := 0
	count := min(x, p.w-x)
	for i := 0; i < count; i++ {
		diffs += p.DiffCols(x-i-1, x+i)
	}
	return diffs
}

func (p puz) DiffsFoldBeforeRow(y int) int {
	// y represents the col after the fold, so y=2, h big: folds 1 onto 2
	// y=7, h=10: 6..7 folds onto 8..9
	diffs := 0
	count := min(y, p.h-y)
	for i := 0; i < count; i++ {
		diffs += p.DiffRows(y-i-1, y+i)
	}
	return diffs
}

func (p puz) DiffCols(x1 int, x2 int) int {
	diffs := 0
	for y := 0; y < p.h; y++ {
		if p.at(pos{x1, y}) != p.at(pos{x2, y}) {
			diffs++
		}
	}
	return diffs
}

func (p puz) DiffRows(y1 int, y2 int) int {
	diffs := 0
	for x := 0; x < p.w; x++ {
		if p.at(pos{x, y1}) != p.at(pos{x, y2}) {
			diffs++
		}
	}
	return diffs
}

type pos struct {
	x, y int
}

func (p pos) dist(o pos) int {
	return abs(p.x-o.x) + abs(p.y-o.y)
}

type wat byte

func abs(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
