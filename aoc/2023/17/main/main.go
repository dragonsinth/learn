package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/sliceheap"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"slices"
	"strings"
)

var sample = `
2413432311323
3215453535623
3255245654254
3446585845452
4546657867536
1438598798454
4457876987766
3637877979653
4654967986887
4564679986453
1224686865563
2546548887735
4322674655533
`

var sample2 = `
111111111111
999999999991
999999999991
999999999991
999999999991
`

func main() {
	run(sample, 1, 3, true)
	run(sample, 4, 10, true)
	run(sample2, 4, 10, true)
}

func run(input string, mn int, mx int, debug bool) {
	p := parse(input)
	if debug {
		termbox.RenderPlain(p.data, os.Stdout)
		fmt.Println()
	}
	start1 := task{0, 0, E, 0, 0, nil}
	start2 := task{0, 0, S, 0, 0, nil}
	t := p.walk([]task{start1, start2}, mn, mx)
	fmt.Println(t.score)

	if debug {
		termbox.RenderPlain(p.renderWith(t), os.Stdout)
		fmt.Println()
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

func (p puz) walk(start []task, mn int, mx int) task {
	work := sliceheap.New(func(a, b task) bool {
		if a.score != b.score {
			return a.score < b.score
		}
		if a.x+a.y != b.x+b.y {
			return a.x+a.y > b.x+b.y
		}
		return false
	})
	for _, t := range start {
		work.Push(t)
	}

	tryWalk := func(t task) {
		if t.x < 0 || t.x >= p.w {
			return
		}
		if t.y < 0 || t.y >= p.h {
			return
		}
		t.score += int(p.data[t.y][t.x] - '0')
		work.Push(t)
	}

	seen := map[taskKey]int{}
	for work.Len() > 0 {
		t := work.Pop()
		if v, ok := seen[t.key()]; ok {
			if v > t.score {
				panic("unexpected")
			}
			continue
		}
		if t.x == p.w-1 && t.y == p.h-1 && t.line >= mn && t.line < mx {
			return t
		}
		seen[t.key()] = t.score

		// Attempt left, right, straight
		if t.line >= mn {
			tryWalk(t.left())
			tryWalk(t.right())
		}
		if t.line < mx {
			tryWalk(t.next())
		}
	}

	panic("unexpected")
}

func (p puz) renderWith(t task) [][]byte {
	out := slices.Clone(p.data)
	for i := range out {
		out[i] = slices.Clone(out[i])
	}
	for _, tk := range t.route {
		out[tk.y][tk.x] = "^>v<"[tk.d]
	}
	return out
}

type task struct {
	x, y  int
	d     dir
	line  int
	score int

	route []taskKey
}

func (t task) next() task {
	switch t.d {
	case N:
		t.y--
		t.line++
	case E:
		t.x++
		t.line++
	case S:
		t.y++
		t.line++
	case W:
		t.x--
		t.line++
	default:
		panic(t.d)
	}
	route := append([]taskKey{}, t.route...)
	t.route = append(route, t.key())
	return t
}

func (t task) left() task {
	t.d = t.d.left()
	t.line = 0
	return t.next()
}

func (t task) right() task {
	t.d = t.d.right()
	t.line = 0
	return t.next()
}

type taskKey struct {
	x, y int
	d    dir
	line int
}

func (t task) key() taskKey {
	return taskKey{
		x:    t.x,
		y:    t.y,
		d:    t.d,
		line: t.line,
	}
}

type dir byte

func (d dir) left() dir {
	return (d + 4 - 1) % 4
}

func (d dir) right() dir {
	return (d + 1) % 4
}

const (
	N = dir(iota)
	E
	S
	W
)
