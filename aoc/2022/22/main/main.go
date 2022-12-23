package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"strconv"
	"strings"
	"time"
)

type input struct {
	width      int // in 2d space
	rsize      int // 3d face size
	regions    [6]region
	layout     string
	directions string
}

var sample = input{
	width: 16,
	rsize: 4,
	regions: [6]region{
		{id: 0, x: 8, y: 0, translate: [4]translation{{5, WEST}, {3, SOUTH}, {2, SOUTH}, {1, SOUTH}}},
		{id: 1, x: 0, y: 4, translate: [4]translation{{2, EAST}, {4, NORTH}, {5, NORTH}, {0, SOUTH}}},
		{id: 2, x: 4, y: 4, translate: [4]translation{{3, EAST}, {4, EAST}, {1, WEST}, {0, EAST}}},
		{id: 3, x: 8, y: 4, translate: [4]translation{{5, SOUTH}, {4, SOUTH}, {2, WEST}, {0, NORTH}}},
		{id: 4, x: 8, y: 8, translate: [4]translation{{5, EAST}, {1, NORTH}, {2, NORTH}, {3, NORTH}}},
		{id: 5, x: 12, y: 8, translate: [4]translation{{0, WEST}, {1, EAST}, {4, WEST}, {3, WEST}}},
	},
	layout: `
        ...#
        .#..
        #...
        ....
...#.......#
........#...
..#....#....
..........#.
        ...#....
        .....#..
        .#......
        ......#.
`,
	directions: `10R5L5R10L4R5L5`,
}

func main() {
	runFlat(sample, 100*time.Millisecond)
	//runFlat(data, time.Millisecond)
	runCube(sample, 100*time.Millisecond)
	//runCube(data, time.Millisecond)
}

func runFlat(in input, frameDelay time.Duration) {
	var p flat
	for _, line := range strings.Split(in.layout, "\n") {
		if line == "" {
			continue
		}

		if len(line) < in.width {
			line += strings.Repeat(" ", in.width-len(line))
		} else if len(line) > in.width {
			panic(line)
		}
		p.data = append(p.data, []byte(line))
	}

	p.width = in.width
	p.height = len(p.data)
	p.computeBounds()

	player := flatPlayer{
		pos: pos{
			x: p.xMins[0],
			y: 0,
		},
		d: EAST,
	}
	p.mark(player)

	tb := termbox.New(true)
	tb.Render(p.data, os.Stdout)

	l := lex{
		in:  []byte(in.directions),
		pos: 0,
	}

	for {
		dist := l.nextInt()
		player = p.walk(player, dist, func() {
			tb.Render(p.data, os.Stdout)
			time.Sleep(frameDelay)
		})

		t := l.nextDir()
		if t == 0 {
			break
		}
		player.d = player.d.turn(t)
		p.mark(player)
	}
	tb.Stop()

	row, col := player.y+1, player.x+1
	fmt.Println(row, col, player.d.score(), 1000*row+4*col+player.d.score())
}

func runCube(in input, frameDelay time.Duration) {
	var c cube
	for _, line := range strings.Split(in.layout, "\n") {
		if line == "" {
			continue
		}

		if len(line) < in.width {
			line += strings.Repeat(" ", in.width-len(line))
		} else if len(line) > in.width {
			panic(line)
		}
		c.data = append(c.data, []byte(line))
	}

	c.rsize = in.rsize
	c.regions = in.regions
	termbox.RenderPlain(c.checkRegions(), os.Stdout)

	player := cubePlayer{
		cp: cubePos{
			pos:      pos{0, 0},
			regionId: 0,
		},
		d: EAST,
	}

	c.mark(player)
	tb := termbox.New(true)
	tb.Render(c.data, os.Stdout)

	l := lex{
		in:  []byte(in.directions),
		pos: 0,
	}

	for {
		dist := l.nextInt()
		player = c.walk(player, dist, func() {
			tb.Render(c.data, os.Stdout)
			time.Sleep(frameDelay)
		})

		t := l.nextDir()
		if t == 0 {
			break
		}
		player.d = player.d.turn(t)
		c.mark(player)
	}
	tb.Stop()

	termbox.RenderPlain(c.data, os.Stdout)

	pt := c.to2d(player.cp)
	row, col := pt.y+1, pt.x+1
	fmt.Println(row, col, player.d.score(), 1000*row+4*col+player.d.score())
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
