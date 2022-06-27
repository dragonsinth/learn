package main

import (
	"fmt"
	"strings"
)

func (p *puzzle) Solve(edgeCounts map[uint16]int) *sol {
	arranged := make([]*tile, p.size*p.size)
	orients := make([]orient, p.size*p.size)
	edges := make([][4]uint16, p.size*p.size)
	for y := 0; y < p.size; y++ {
		for x := 0; x < p.size; x++ {
			if x == 1 && y == 1 {
				fmt.Println("here")
			}

			var f func(t *tile, or orient) bool
			if x == 0 && y == 0 {
				// Must orient the tile so that left and top are uniq
				f = func(t *tile, or orient) bool {
					if t.typ != CORNER {
						return false
					}
					if or <= WEST {
						return false
					}
					if edgeCounts[t.edges[or][NORT]] == 1 &&
						edgeCounts[t.edges[or][WEST]] == 1 {
						return true
					}
					return false
				}
			} else if x > 0 {
				// find the tile that matches to the left
				want := edges[y*p.size+x-1][EAST]
				want = invert(want)
				f = func(t *tile, or orient) bool {
					if t.edges[or][WEST] == want {
						return true
					}
					return false
				}
			} else if y > 0 {
				// find the tile that matches above
				want := edges[y*p.size+x-p.size][SOUT]
				want = invert(want)
				f = func(t *tile, or orient) bool {
					if t.edges[or][NORT] == want {
						return true
					}
					return false
				}
			}

			t, or := p.find(f)
			fmt.Printf("(%d,%d) %d %s %s, %+v\n", y, x, t.id, t.typ.String(), or.String(), t.edges[or])
			arranged[y*p.size+x] = t
			orients[y*p.size+x] = or
			edges[y*p.size+x] = t.edges[or]
		}
	}

	var s sol
	for y := 0; y < p.size; y++ {
		var bits [][][]bool
		for x := 0; x < p.size; x++ {
			bits = append(bits, arranged[y*p.size+x].Render(orients[y*p.size+x]))
		}

		for ly := 0; ly < isize; ly++ {
			var line []bool
			for _, bit := range bits {
				line = append(line, bit[ly][:]...)
			}
			s.lines = append(s.lines, line)
		}
	}

	return &s
}

func (p *puzzle) find(f func(t *tile, or orient) bool) (*tile, orient) {
	for i, t := range p.tiles {
		for o := NORT; o < MAX_ORIENT; o++ {
			if f(t, o) {
				p.tiles = append(p.tiles[:i], p.tiles[i+1:]...)
				return t, o
			}
		}
	}
	panic("not found")
}

type sol struct {
	lines [][]bool
}

func (s *sol) Solve() int {
	monsters := [MAX_ORIENT][][]bool{}
	for or := NORT; or < MAX_ORIENT; or++ {
		monsters[or] = rotate(monster, or)
	}

	// Copy the entire input.
	cp := make([][]bool, len(s.lines))
	for i, line := range s.lines {
		cp[i] = append([]bool{}, line...)
	}

	ymax, xmax := len(s.lines), len(s.lines[0])

	for y := 0; y < ymax; y++ {
		for x := 0; x < xmax; x++ {
			// See if any monster variant matches
			for or := NORT; or < MAX_ORIENT; or++ {
				m := monsters[or]
				if y+len(m) > ymax || x+len(m[0]) > xmax {
					// Does not fit
					continue
				}

				match := true
				for j := 0; j < len(m); j++ {
					for i := 0; i < len(m[0]); i++ {
						if m[j][i] && !s.lines[y+j][x+i] {
							match = false
						}
					}
				}
				if match {
					fmt.Println(x, y, or, "match")

					// Blot out the copy (but not the original)
					for j := 0; j < len(m); j++ {
						for i := 0; i < len(m[0]); i++ {
							if m[j][i] {
								cp[y+j][x+i] = false
							}
						}
					}
				}
			}
		}
	}

	// Count the remaining dots
	sum := 0
	for _, line := range cp {
		for _, v := range line {
			if v {
				sum++
			}
		}
	}
	return sum
}

func (s *sol) String() string {
	return toString(s.lines)
}

var monster = initMonster(`
..................#.
#....##....##....###
.#..#..#..#..#..#...
`)

func initMonster(txt string) [][]bool {
	ret := makeBoolField(3, 20)
	wy := 0
	for _, line := range strings.Split(txt, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		for i, c := range line {
			if c == '#' {
				ret[wy][i] = true
			} else {
				ret[wy][i] = false
			}
		}
		wy++
	}
	return ret
}

func makeBoolField(y int, x int) [][]bool {
	ret := make([][]bool, y)
	for i := 0; i < y; i++ {
		ret[i] = make([]bool, x)
	}
	return ret
}
