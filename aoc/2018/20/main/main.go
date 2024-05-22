package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
)

var samples = []struct {
	expect int
	input  string
}{
	{3, `^WNE$`},
	{10, `^ENWWW(NEEE|SSE(EE|N))$`},
	{18, `^ENNWSWW(NEWS|)SSSEEN(WNSE|)EE(SWEN|)NNN$`},
	{23, `^ESSWWN(E|NNENN(EESS(WNSE|)SSS|WWWSSSSE(SW|NNNE)))$`},
	{31, `^WSSEESWWWNW(S|NENNEEEENN(ESSSSW(NWSW|SSEN)|WSWWN(E|WWS(E|SS))))$`},
}

func main() {
	for _, sample := range samples {
		distances := run(sample.input, true)
		fmt.Println(farthest(distances), atLeast1000(distances))
	}
}

func run(s string, debug bool) map[pos]int {
	p := puzzle{
		wat: map[pos]room{},
		min: pos{x: 0, y: 0},
		max: pos{x: 0, y: 0},
	}

	in := &input{s, 0}
	if c := in.Peek(); c != '^' {
		panic(c)
	}
	in.Next()

	origin := pos{0, 0}
	p.Walk(in, origin)
	if c := in.Peek(); c != '$' {
		panic(c)
	}
	in.Next()
	if !in.Empty() {
		panic(in.Peek())
	}

	if debug {
		termbox.RenderPlain(p.Render(), os.Stdout)
	}

	return p.distances(origin)
}

type input struct {
	s string
	p int
}

func (in *input) Next() {
	if in.Empty() {
		panic("empty")
	}
	in.p++
}

func (in *input) Peek() byte {
	return in.s[in.p]
}

func (in *input) Empty() bool {
	return in.p >= len(in.s)
}

type dir byte

const (
	N = dir(iota)
	E
	S
	W
)

type room byte

func (r room) With(d dir) room {
	return r | (1 << d)
}

func (r room) Has(d dir) bool {
	return (r & (1 << d)) != 0
}

func (r room) String() string {
	return "NSEW"[r : r+1]
}

type pos struct {
	x, y int
}

func (p pos) Walk(d dir) pos {
	switch d {
	case N:
		return pos{p.x, p.y - 1}
	case E:
		return pos{p.x + 1, p.y}
	case S:
		return pos{p.x, p.y + 1}
	case W:
		return pos{p.x - 1, p.y}
	default:
		panic(d)
	}
}

func (p pos) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}
