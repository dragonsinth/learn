package main

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"time"
)

type puz struct {
	data [][]byte
	w, h int
	turn int

	term termbox.Terminal
}

func (p *puz) Tick() {
	out := make([][]byte, p.h)
	for y := 0; y < p.h; y++ {
		out[y] = make([]byte, p.w)
		for x := 0; x < p.w; x++ {
			nTree, nLumber := p.getAdjacent(pos{x, y})
			switch p.data[y][x] {
			case '.':
				if nTree >= 3 {
					out[y][x] = '|'
				} else {
					out[y][x] = '.'
				}
			case '|':
				if nLumber >= 3 {
					out[y][x] = '#'
				} else {
					out[y][x] = '|'
				}
			case '#':
				if nTree >= 1 && nLumber >= 1 {
					out[y][x] = '#'
				} else {
					out[y][x] = '.'
				}
			default:
				panic(p.data[y][x])
			}
		}
	}
	p.data = out
}

func (p *puz) Print() {
	p.term.Render(p.Render(), os.Stdout)
}

func (p *puz) PrintFrame() {
	p.term.Render(p.Render(), os.Stdout)
	if p.term.Enabled() {
		p.wait()
	}
}

func (p *puz) Render() [][]byte {
	return append(p.data, nil)
}

func (p *puz) Count() (nTree, nLumber int) {
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			switch p.data[y][x] {
			case '.':
			case '|':
				nTree++
			case '#':
				nLumber++
			default:
				panic(p.data[y][x])
			}
		}
	}
	return
}

func (p *puz) wait() {
	time.Sleep(100 * time.Millisecond)
}

func (p *puz) getAdjacent(from pos) (nTree, nLumber int) {
	for _, pt := range from.adjacent() {
		if pt.x < 0 || pt.x >= p.w {
			continue
		}
		if pt.y < 0 || pt.y >= p.h {
			continue
		}
		switch p.data[pt.y][pt.x] {
		case '.':
		case '|':
			nTree++
		case '#':
			nLumber++
		default:
			panic(p.data[pt.y][pt.x])
		}
	}
	return
}

func (p *puz) key() string {
	h := sha256.New()
	for _, line := range p.data {
		h.Write(line)
	}
	out := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(out)
}
