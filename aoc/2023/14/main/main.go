package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"strings"
)

var sample = `
O....#....
O.OO#....#
.....##...
OO.#O....O
.O.....O#.
O.#..O.#.#
..O..#O..O
.......O..
#....###..
#OO..#....
`

func main() {
	part1(sample, true)
	part2(sample, true)
}

func part1(input string, debug bool) {
	p := parse(input)
	if debug {
		termbox.RenderPlain(p.Render(), os.Stdout)
		fmt.Println()
	}
	p.Tilt(0, -1)
	if debug {
		fmt.Println("TiltNorth")
		termbox.RenderPlain(p.Render(), os.Stdout)
		fmt.Println()
	}
	fmt.Println(p.CountLoad())
	fmt.Println()
}

func part2(input string, debug bool) {
	var vals []int
	seen := map[string]int{}
	p := parse(input)
	var i int
	for {
		sha := p.Sha()
		load := p.CountLoad()
		if debug {
			fmt.Println(i, sha, load)
		}
		vals = append(vals, load)
		last := seen[sha]
		if last != 0 {
			fmt.Printf("%s=%d, previous=%d, now=%d\n", sha, load, last, i)
			mod := i - last
			where := (1000000000-last)%mod + last
			fmt.Printf("mod=%d, where=%d, val=%d\n", mod, where, vals[where])
			return
		}
		seen[sha] = i

		p.Tilt(0, -1)
		p.Tilt(-1, 0)
		p.Tilt(0, 1)
		p.Tilt(1, 0)
		i++
	}
}

func parse(input string) puz {
	var p puz
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if p.w == 0 {
			p.w = len(line)
		} else if p.w != len(line) {
			panic(line)
		}

		p.data = append(p.data, []byte(line))
		p.h++
	}
	return p
}

type puz struct {
	data [][]byte
	w, h int
}

func (p puz) Render() [][]byte {
	return p.data
}

func (p puz) Tilt(dx, dy int) {
	xs, xd := 0, 1
	ys, yd := 0, 1
	if dx > 0 {
		// to roll east, we iterate west -> east
		xs, xd = p.w-1, -1
	}
	if dy > 0 {
		// to roll south, we iterate south -> north
		ys, yd = p.h-1, -1
	}

	changed := true
	for changed {
		changed = false

		for y := ys; y >= 0 && y < p.h; y += yd {
			for x := xs; x >= 0 && x < p.w; x += xd {
				tx, ty := x+dx, y+dy
				if tx < 0 || tx >= p.w {
					continue
				}
				if ty < 0 || ty >= p.h {
					continue
				}

				if p.data[y][x] == 'O' {
					// see if it can move north
					if p.data[ty][tx] == '.' {
						changed = true
						p.data[ty][tx] = 'O'
						p.data[y][x] = '.'
					}
				}
			}
		}
	}
}

func (p puz) CountLoad() int {
	sum := 0
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			if p.data[y][x] == 'O' {
				sum += p.h - y
			}
		}
	}
	return sum
}

func (p puz) Sha() string {
	bufs := p.Render()
	h := sha256.New()
	for _, b := range bufs {
		h.Write(b)
	}
	return hex.EncodeToString(h.Sum(nil))
}
