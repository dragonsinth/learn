package main

import (
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"strconv"
	"strings"
)

func part2(input string, w, h int, turns int, debug bool) {
	p := puz{
		data: nil,
		w:    w,
		h:    h,
	}
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if !re.MatchString(line) {
			panic(line)
		}
		matches := re.FindStringSubmatch(line)
		if len(matches) != 5 {
			panic(matches)
		}

		var r robot
		r.x, r.y, r.dx, r.dy = mustInt(matches[1]), mustInt(matches[2]), mustInt(matches[3]), mustInt(matches[4])
		p.data = append(p.data, r)
	}

	p.Run(turns, debug)
}

type robot struct {
	x, y   int
	dx, dy int
}

type puz struct {
	data []robot
	w, h int
	turn int
}

func (p *puz) Tick() {
	for i := range p.data {
		r := &p.data[i]
		r.x = mod(r.x+r.dx, p.w)
		r.y = mod(r.y+r.dy, p.h)
	}
}

type pos struct {
	x, y int
}

func (p *puz) Render() [][]byte {
	rpos := make(map[pos]bool, len(p.data))
	for _, p := range p.data {
		rpos[pos{p.x, p.y}] = true
	}
	var buf [][]byte
	for y := 0; y < p.h; y++ {
		var line []byte
		for x := 0; x < p.w; x++ {
			pt := pos{x, y}
			if rpos[pt] {
				line = append(line, '#')
			} else {
				line = append(line, ' ')
			}
		}
		buf = append(buf, line)
	}
	buf = append(buf, []byte(strconv.Itoa(p.turn)))

	return buf
}

func (p *puz) Run(turns int, debug bool) {
	term := termbox.New(debug)
	defer func() {
		term.Stop()
		term.Render(p.Render(), os.Stdout)
	}()
	defer term.Stop()
	for p.turn < turns {
		if debug {
			term.Render(p.Render(), os.Stdout)
			b := [1]byte{}
			_, _ = os.Stdin.Read(b[:])
			if b[0] == 'q' {
				break
			}
		}
		p.Tick()
		p.turn++
	}
}
