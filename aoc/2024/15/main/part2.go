package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"strings"
)

func part2(input string, debug bool) {
	p, instr := parse2(input)
	p.Run(instr, debug)
}

func parse2(input string) (*puz2, string) {
	p := &puz2{
		data: map[pos]byte{},
	}
	instr := ""

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "#") {
			if p.w <= 0 {
				p.w = len(line) * 2
			} else if p.w != len(line)*2 {
				panic(line)
			}
			y := p.h
			for x, wat := range line {
				pta := pos{x: 2 * x, y: y}
				ptb := pos{x: 2*x + 1, y: y}
				switch wat {
				case '#':
					p.data[pta] = '#'
					p.data[ptb] = '#'
				case 'O':
					p.data[pta] = '['
					p.data[ptb] = ']'
				case '.':
					p.data[pta] = '.'
					p.data[ptb] = '.'
				case '@':
					p.bot = pta
					p.data[pta] = '@'
					p.data[ptb] = '.'
				default:
					panic(wat)
				}
			}
			p.h++
		} else {
			instr += line
		}

	}
	return p, instr
}

type puz2 struct {
	data map[pos]byte
	w, h int
	bot  pos
}

func (p *puz2) Run(instr string, debug bool) {
	term := termbox.New(debug)
	defer func() {
		term.Stop()
		term.Render(p.Render(), os.Stdout)
		fmt.Println(p.Score())
	}()
	defer term.Stop()
	for _, d := range instr {
		if debug {
			term.Render(p.Render(), os.Stdout)
			b := [1]byte{}
			_, _ = os.Stdin.Read(b[:])
			if b[0] == 'q' {
				break
			}
		}
		p.Tick(dir(d))
	}
}

func (p *puz2) Tick(d dir) {
	if p.CanMove(p.bot, d) {
		p.Move(p.bot, d)
		dst := p.bot.next(d)
		assert(p.data[dst] == '@')
		p.bot = dst
	}
}

func (p *puz2) CanMove(src pos, d dir) bool {
	dst := src.next(d)
	switch p.data[dst] {
	case '.':
		// It's empty, do the move
		return true
	case '#':
		return false
	case '[':
		// check the left side of the box first
		if !p.CanMove(dst, d) {
			return false
		}
		// if pushing north or south, also check the right half of the box
		if d == N || d == S {
			right := pos{dst.x + 1, dst.y}
			if !p.CanMove(right, d) {
				return false
			}
		}
		return true
	case ']':
		// check the right side of the box first
		if !p.CanMove(dst, d) {
			return false
		}
		// if pushing north or south, also check the left half of the box too
		if d == N || d == S {
			left := pos{dst.x - 1, dst.y}
			if !p.CanMove(left, d) {
				return false
			}
		}
		return true
	default:
		panic(string(p.data[dst]))
	}
}

func (p *puz2) Move(src pos, d dir) {
	dst := src.next(d)
	switch p.data[dst] {
	case '.':
		// ok
	case '#':
		panic("wall")
	case '[':
		// recursively move the box
		p.Move(dst, d)
		assert(p.data[dst] == '.')
		// if pushing north or south, also push the right half of the box
		if d == N || d == S {
			right := pos{dst.x + 1, dst.y}
			p.Move(right, d)
			assert(p.data[right] == '.')
		}
	case ']':
		// recursively move the box
		p.Move(dst, d)
		assert(p.data[dst] == '.')
		// if pushing north or south, also push the left half of the box
		if d == N || d == S {
			left := pos{dst.x - 1, dst.y}
			p.Move(left, d)
			assert(p.data[left] == '.')
		}
	default:
		panic(string(p.data[dst]))
	}

	p.data[dst] = p.data[src]
	p.data[src] = '.'
}

func (p *puz2) Render() [][]byte {
	var buf [][]byte
	for y := 0; y < p.h; y++ {
		var line []byte
		for x := 0; x < p.w; x++ {
			line = append(line, p.data[pos{x, y}])
		}
		buf = append(buf, line)
	}
	return buf
}

func (p *puz2) Score() int {
	sum := 0
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			if p.data[pos{x, y}] == '[' {
				sum += x + 100*y
			}
		}
	}
	return sum
}
