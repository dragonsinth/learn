package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"strings"
)

func part1(input string, debug bool) {
	p, instr := parse1(input)
	p.Run(instr, debug)
}

func parse1(input string) (*puz1, string) {
	p := &puz1{
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
				p.w = len(line)
			} else if p.w != len(line) {
				panic(line)
			}
			y := p.h
			for x, wat := range line {
				pt := pos{x: x, y: y}
				p.data[pt] = byte(wat)
				if wat == '@' {
					p.bot = pt
				}
			}
			p.h++
		} else {
			instr += line
		}

	}
	return p, instr
}

type puz1 struct {
	data map[pos]byte
	w, h int
	bot  pos
}

func (p *puz1) Run(instr string, debug bool) {
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

func (p *puz1) Tick(d dir) {
	if p.Move(p.bot, d) {
		dst := p.bot.next(d)
		assert(p.data[dst] == '@')
		p.bot = dst
	}
}

func (p *puz1) Move(src pos, d dir) bool {
	dst := src.next(d)
	switch p.data[dst] {
	case '.':
		// It's empty, do the move
		p.data[dst] = p.data[src]
		p.data[src] = '.'
		return true
	case '#':
		return false
	case 'O':
		// recursively try to move the box
		if p.Move(dst, d) {
			assert(p.data[dst] == '.')
			p.data[dst] = p.data[src]
			p.data[src] = '.'
			return true
		}
		return false
	default:
		panic(string(p.data[dst]))
	}
}

func (p *puz1) Render() [][]byte {
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

func (p *puz1) Score() int {
	sum := 0
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			if p.data[pos{x, y}] == 'O' {
				sum += x + 100*y
			}
		}
	}
	return sum
}
