package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/grid"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var sample = `
1,0,1~1,2,1
0,0,2~2,0,2
0,2,3~2,2,3
0,0,4~0,2,4
2,0,5~2,2,5
0,1,6~2,1,6
1,1,8~1,1,9
`

var (
	re = regexp.MustCompile(`^(\d+),(\d+),(\d+)~(\d+),(\d+),(\d+)$`)
)

func main() {
	run(sample, true)
}

func run(input string, debug bool) {
	p := parse(input)
	if debug {
		termbox.RenderPlain(p.Render(Y, false), os.Stdout)
		termbox.RenderPlain(p.Render(X, false), os.Stdout)
	}

	p.drop()
	if debug {
		termbox.RenderPlain(p.Render(Y, false), os.Stdout)
		termbox.RenderPlain(p.Render(X, false), os.Stdout)
	}

	deps, rdeps := make([][]int, len(p.blocks)), make([][]int, len(p.blocks))
	cannot := make([]bool, len(p.blocks))
	for _, b := range p.blocks {
		if b.min(Z) == 1 {
			continue
		}
		collisions := p.checkCollisions(b, b.min(Z)-1)
		for _, c := range collisions {
			if debug {
				fmt.Printf("%d -> %d;\n", b.id, c.id)
			}
			deps[b.id] = append(deps[b.id], c.id)
			rdeps[c.id] = append(rdeps[c.id], b.id)
		}
		if len(collisions) == 1 {
			supporter := collisions[0]
			cannot[supporter.id] = true
		}
	}
	sum := 0
	for _, v := range cannot {
		if !v {
			sum++
		}
	}
	fmt.Println(sum)

	fmt.Println(graphWalk(deps, rdeps))
}

func parse(input string) puz {
	var p puz
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !re.MatchString(line) {
			panic(line)
		}
		parts := re.FindStringSubmatch(line)
		b := block{
			id: len(p.blocks),
			a: point{
				mustInt(parts[1]),
				mustInt(parts[2]),
				mustInt(parts[3]),
			},
			b: point{
				mustInt(parts[4]),
				mustInt(parts[5]),
				mustInt(parts[6]),
			},
		}
		for d := X; d <= Z; d++ {
			if b.a[d] < b.b[d] {
				b.dim = d
				break
			} else if b.a[d] > b.b[d] {
				b.dim = d
				b.a, b.b = b.b, b.a
				break
			}
		}
		b.len = b.b[b.dim] - b.a[b.dim]
		p.blocks = append(p.blocks, b)
	}
	slices.SortFunc(p.blocks, func(a, b block) int {
		return a.min(Z) - b.min(Z)
	})
	p.bbz = newBlocksByZ(p.blocks)
	return p
}

type puz struct {
	blocks []block
	bbz    blocksByZ
}

func (p puz) Render(d dim, neg bool) [][]byte {
	// sort along the view dimension, from back to front (painters)
	blocks := slices.Clone(p.blocks)
	slices.SortFunc(blocks, func(a, b block) int {
		if neg {
			return a.min(d) - b.min(d)
		} else {
			return b.max(d) - a.max(d)
		}
	})

	var dx, dy dim
	switch d {
	case X:
		dx, dy = Y, Z
	case Y:
		dx, dy = X, Z
	case Z:
		dx, dy = X, Y
	}

	// compute the max in each dimension
	var maxX, maxY int
	for _, b := range blocks {
		maxX = max(maxX, b.max(dx))
		maxY = max(maxY, b.max(dy))
	}
	maxX++
	maxY++

	fmt.Println(dx, maxX, dy, maxY)

	buf := grid.Alloc2d(maxX, maxY, byte('.'))
	for x := 0; x < maxX; x++ {
		buf[maxY-1][x] = '-'
	}
	// paint front to back
	for _, b := range blocks {
		for y := b.min(dy); y <= b.max(dy); y++ {
			for x := b.min(dx); x <= b.max(dx); x++ {
				buf[maxY-y-1][x] = b.tok()
			}
		}
	}

	return buf
}

func (p puz) drop() {
	for i := range p.blocks {
		b := p.blocks[i]
		newBlock, changed := p.dropBlock(b)
		if changed {
			p.bbz.rem(b)
			p.bbz.add(newBlock)
			p.blocks[i] = newBlock
		}
	}
}

func (p puz) dropBlock(b block) (block, bool) {
	changed := false
	for b.min(Z) > 1 {
		collisions := p.checkCollisions(b, b.min(Z)-1)
		if len(collisions) > 0 {
			return b, changed
		}
		b.a[Z]--
		b.b[Z]--
		changed = true
	}
	return b, changed
}

// returns the set of collisions that would occur from dropping block to z.
func (p puz) checkCollisions(b block, z int) []block {
	var ret []block
	m := p.bbz[z]
	for other := range m {
		if b.id == other.id {
			panic("here")
		}
		if b.overlap(other, X) && b.overlap(other, Y) {
			ret = append(ret, other)
		}
	}
	return ret
}

type dim int

const (
	X = dim(0)
	Y = dim(1)
	Z = dim(2)
)

func (d dim) String() string {
	return "XYZ"[d : d+1]
}

type block struct {
	id   int
	a, b point
	dim  dim
	len  int
}

func (b block) tok() byte {
	return byte(b.id) + 'A'
}

func (b block) String() string {
	return fmt.Sprintf("%d: %s(%d); %s->%s", b.id, b.dim, b.len, b.a, b.b)
}

func (b block) min(d dim) int {
	return b.a[d]
}

func (b block) max(d dim) int {
	return b.b[d]
}

func (b block) overlap(other block, d dim) bool {
	return max(b.min(d), other.min(d)) <= min(b.max(d), other.max(d))
}

func (b block) overlapPoint(p int, d dim) bool {
	return max(b.min(d), p) <= min(b.max(d), p)
}

type point [3]int

func (p point) String() string {
	return fmt.Sprintf("%d,%d,%d", p[X], p[Y], p[Z])
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
