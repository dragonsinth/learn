package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"os"
	"strings"
	"time"
)

const sample = `
/->-\        
|   |  /----\
| /-+--+-\  |
| | |  | v  |
\-+-/  \-+--/
  \------/   
`

const sample2 = `
/>-<\  
|   |  
| /<+-\
| | | v
\>+</ |
  |   ^
  \<->/
`

func main() {
	// part 1
	p := parse(sample)
	term := termbox.New(true)
	term.Render(p.Render(), os.Stdout)
	for i := 1; true; i++ {
		p.Tick()
		if i < 20 {
			time.Sleep(100 * time.Millisecond)
			term.Render(p.Render(), os.Stdout)
		}
		if p.crash != nil {
			term.Stop()
			term.Render(p.Render(), os.Stdout)
			fmt.Println(i, p.crash)
			break
		}
	}

	// part 2
	p = parse(sample2)
	term.Start()
	term.Render(p.Render(), os.Stdout)
	for i := 1; true; i++ {
		p.Tick()
		if i < 20 {
			time.Sleep(100 * time.Millisecond)
			term.Render(p.Render(), os.Stdout)
		}
		if len(p.cars) < 2 {
			term.Stop()
			term.Render(p.Render(), os.Stdout)
			for _, c := range p.cars {
				fmt.Println(i, p.crash, c.p)
			}
			break
		}
	}
}

func parse(input string) puz {
	p := puz{
		cars: map[pos]car{},
		road: map[pos]byte{},
		w:    -1,
		h:    -1,
	}

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if p.w < 0 {
			p.w = len(line)
		} else if p.w != len(line) {
			panic(line)
		}
		p.h++

		for x, c := range []byte(line) {
			pt := pos{x, p.h}
			switch c {
			case '/', '\\', '|', '-', '+':
				p.road[pt] = c
			case '^':
				p.road[pt] = '|'
				p.cars[pt] = car{p: pt, d: NORTH}
			case '>':
				p.road[pt] = '-'
				p.cars[pt] = car{p: pt, d: EAST}
			case 'v':
				p.road[pt] = '|'
				p.cars[pt] = car{p: pt, d: SOUTH}
			case '<':
				p.road[pt] = '-'
				p.cars[pt] = car{p: pt, d: WEST}
			case ' ':
			default:
				panic(string(c))
			}
		}
	}
	p.h++
	return p
}

type pos struct {
	x, y int
}

type dir byte

const (
	NORTH = dir(iota)
	EAST
	SOUTH
	WEST
)

func (d dir) Char() byte {
	switch d {
	case NORTH:
		return '^'
	case EAST:
		return '>'
	case SOUTH:
		return 'v'
	case WEST:
		return '<'
	default:
		panic(d)
	}
}

func (d dir) Left() dir {
	return (d + 3) % 4
}

func (d dir) Right() dir {
	return (d + 1) % 4
}

type turn byte

const (
	LEFT = turn(iota)
	STRAIGHT
	RIGHT
)

type car struct {
	p pos
	d dir
	t turn
}

func (c car) Next() pos {
	p := c.p
	switch c.d {
	case NORTH:
		p.y--
	case EAST:
		p.x++
	case SOUTH:
		p.y++
	case WEST:
		p.x--
	default:
		panic(c.d)
	}
	return p
}

func carSort(a, b car) bool {
	if a.p.y != b.p.y {
		return a.p.y < b.p.y
	}
	return a.p.x < b.p.x
}

func sortCars(cars []car) []car {
	slices.SortFunc(cars, carSort)
	return cars
}

type puz struct {
	cars  map[pos]car
	road  map[pos]byte
	crash *pos
	w, h  int
}

func (p *puz) Tick() {
	// run each car in order
	for _, c := range sortCars(maps.Values(p.cars)) {
		if _, ok := p.cars[c.p]; !ok {
			// must have crashed
			continue
		}

		// move
		delete(p.cars, c.p)
		dst := c.Next()

		if _, ok := p.cars[dst]; ok {
			// crash!
			p.crash = &dst

			// remove the other car as well
			delete(p.cars, dst)
			continue
		}

		c.p = dst

		// compute next direction
		switch p.road[dst] {
		case '-', '|':
		// nothing
		case '/':
			switch c.d {
			case EAST, WEST:
				c.d = c.d.Left()
			case NORTH, SOUTH:
				c.d = c.d.Right()
			}
		case '\\':
			switch c.d {
			case EAST, WEST:
				c.d = c.d.Right()
			case NORTH, SOUTH:
				c.d = c.d.Left()
			}
		case '+':
			switch c.t {
			case LEFT:
				c.d = c.d.Left()
			case STRAIGHT:
			case RIGHT:
				c.d = c.d.Right()
			default:
				panic(c.t)
			}
			c.t = (c.t + 1) % 3
		default:
			panic(p.road[dst])
		}

		p.cars[dst] = c
	}
}

func (p *puz) Render() [][]byte {
	var buf [][]byte
	for y := 0; y < p.h; y++ {
		var line []byte
		for x := 0; x < p.w; x++ {
			var ch byte
			pt := pos{x, y}
			if p.crash != nil && *p.crash == pt {
				ch = 'X'
			} else if c, ok := p.cars[pt]; ok {
				ch = c.d.Char()
			} else if r, ok := p.road[pt]; ok {
				ch = r
			} else {
				ch = ' '
			}
			line = append(line, ch)
		}
		buf = append(buf, line)
	}
	return buf
}
