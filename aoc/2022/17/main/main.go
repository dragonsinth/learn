package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"strings"
	"time"
)

var sample = `>>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>`

type rock struct {
	w, h int
	data [][]byte
}

var rocks = []rock{
	{4, 1, [][]byte{
		[]byte("####"),
	}},
	{3, 3, [][]byte{
		[]byte(".#."),
		[]byte("###"),
		[]byte(".#."),
	}},
	{3, 3, [][]byte{
		[]byte("###"),
		[]byte("..#"),
		[]byte("..#"),
	}},
	{1, 4, [][]byte{
		[]byte("#"),
		[]byte("#"),
		[]byte("#"),
		[]byte("#"),
	}},
	{2, 2, [][]byte{
		[]byte("##"),
		[]byte("##"),
	}},
}

func main() {
	run(sample, 10*time.Millisecond, 12)

	sum, dx, dy := run(sample, 0, 2022)
	fmt.Println(sum)
	toMod, _, _ := run(sample, 0, 1000000000000%dx)
	fmt.Println(toMod)
	fmt.Println((1000000000000/dx)*dy + toMod)
}

func run(input string, dbgFrameTime time.Duration, maxRocks int) (int, int, int) {
	term := termbox.New(dbgFrameTime != 0)
	defer term.Stop()
	defer time.Sleep(100 * dbgFrameTime)

	input = strings.TrimSpace(input)
	p := puzzle{
		input:      input,
		inPos:      0,
		inLen:      len(input),
		nextRock:   0,
		rockHeight: 0,
		data: [][]byte{
			[]byte("+-------+"),
		},
	}

	type val struct {
		turn   int
		height int
	}

	seen := map[key]val{}
	var dx, dy int
	for i := 0; i < maxRocks; i++ {
		k := p.key()
		v := val{
			turn:   i,
			height: p.rockHeight,
		}
		if dx != 0 {
			// nothing else to do
		} else if prev, ok := seen[k]; ok {
			dx = v.turn - prev.turn
			dy = v.height - prev.height
			fmt.Printf("state from %d,%d seen again at %d,%d; dx=%d,dy=%d\n",
				prev.turn, prev.height,
				v.turn, v.height,
				dx, dy,
			)
			seen = nil
		} else {
			seen[k] = v
		}

		p.dropRock(term, dbgFrameTime)
	}
	return p.rockHeight, dx, dy
}

type puzzle struct {
	input      string
	inPos      int
	inLen      int
	nextRock   int
	rockHeight int
	data       [][]byte

	rock    *rock
	rockPos pos
}

func (p *puzzle) dropRock(term termbox.Terminal, dbgFrameTime time.Duration) {
	dbgFrame := func() {
		if dbgFrameTime != 0 {
			term.Render(p.renderView(), os.Stdout)
			time.Sleep(dbgFrameTime)
		}
	}

	p.rock = &rocks[p.nextRock]
	p.nextRock = (p.nextRock + 1) % len(rocks)

	p.rockPos = pos{3, p.rockHeight + 3 + 1}
	p.ensure(p.rockPos.y + p.rock.h)
	dbgFrame()

	for {
		// blow rock left<->right
		if p.input[p.inPos] == '<' {
			// move left
			newPos := p.rockPos
			newPos.x--
			if !p.clip(newPos) {
				p.rockPos = newPos
			}
		} else if p.input[p.inPos] == '>' {
			// move right
			newPos := p.rockPos
			newPos.x++
			if !p.clip(newPos) {
				p.rockPos = newPos
			}
		} else {
			panic(p.input[p.inPos])
		}

		p.inPos = (p.inPos + 1) % p.inLen
		dbgFrame()

		// drop
		newPos := p.rockPos
		newPos.y--
		if !p.clip(newPos) {
			p.rockPos = newPos
			dbgFrame()
		} else {
			p.commit()
			break
		}
	}

	p.rock = nil
	dbgFrame()
}

func (p *puzzle) ensure(y int) {
	for len(p.data) < y {
		p.data = append(p.data, []byte("|.......|"))
	}
}

func (p *puzzle) renderView() [][]byte {
	// Must render backwards.
	sz := len(p.data)
	outRow := func(inRow int) int {
		if inRow < 0 || inRow >= sz {
			panic(inRow)
		}
		return sz - inRow - 1
	}

	out := make([][]byte, sz)
	for i := 0; i < len(out); i++ {
		out[outRow(i)] = p.data[i]
	}

	// TODO: render active block.
	if p.rock != nil {
		for y := 0; y < p.rock.h; y++ {
			// copy-on-write the outrow so we don't destroy the game state
			outY := outRow(y + p.rockPos.y)
			out[outY] = append([]byte{}, out[outY]...)
			for x := 0; x < p.rock.w; x++ {
				outX := x + p.rockPos.x
				if p.rock.data[y][x] == '#' {
					out[outY][outX] = '@'
				}
			}
		}
	}
	return out
}

func (p *puzzle) clip(newPos pos) bool {
	for y := 0; y < p.rock.h; y++ {
		py := y + newPos.y
		for x := 0; x < p.rock.w; x++ {
			px := x + newPos.x
			if p.rock.data[y][x] == '#' {
				if p.data[py][px] != '.' {
					return true // clipped
				}
			}
		}
	}
	return false
}

func (p *puzzle) commit() {
	for y := 0; y < p.rock.h; y++ {
		py := y + p.rockPos.y
		for x := 0; x < p.rock.w; x++ {
			px := x + p.rockPos.x
			if p.rock.data[y][x] == '#' {
				p.data[py][px] = '#'
			}
		}
	}
	p.rockHeight = max(p.rockHeight, p.rockPos.y+p.rock.h-1)
}

type key struct {
	nextRock int
	inPos    int
	topLine  string
}

func (p *puzzle) key() key {
	return key{
		nextRock: p.nextRock,
		inPos:    p.inPos,
		topLine:  string(p.data[p.rockHeight]),
	}
}

type pos struct {
	x, y int
}

func (pt pos) String() string {
	return fmt.Sprintf("(%d,%d)", pt.x, pt.y)
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
