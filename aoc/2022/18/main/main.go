package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"math"
	"os"
	"strconv"
	"strings"
)

var sample = `
2,2,2
1,2,2
3,2,2
2,1,2
2,3,2
2,2,1
2,2,3
2,2,4
2,2,6
1,2,5
3,2,5
2,1,5
2,3,5
`

type wat int

const (
	AIR = wat(iota)
	LAVA
	STEAM
)

func main() {
	p := puzzle{
		data: map[pos]wat{},
		min:  pos{math.MaxInt, math.MaxInt, math.MaxInt},
	}

	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		pt := pos{
			x: mustInt(parts[0]),
			y: mustInt(parts[1]),
			z: mustInt(parts[2]),
		}
		p.data[pt] = LAVA
		p.min.x = min(p.min.x, pt.x)
		p.min.y = min(p.min.y, pt.y)
		p.min.z = min(p.min.z, pt.z)
		p.max.x = max(p.max.x, pt.x)
		p.max.y = max(p.max.y, pt.y)
		p.max.z = max(p.max.z, pt.z)
	}
	p.min.x--
	p.min.y--
	p.min.z--
	p.max.x++
	p.max.y++
	p.max.z++
	fmt.Println(p.min, p.max)

	// Part 1
	fmt.Println(p.countSurface(AIR))

	// Part 2
	// 3d flood-fill
	p.fill(p.min)

	for z := p.min.z; z <= p.max.z; z++ {
		termbox.RenderPlain(p.renderPlane(z), os.Stdout)
		fmt.Println()
	}

	fmt.Println(p.countSurface(STEAM))
}

type pos struct {
	x, y, z int
}

func (pt pos) String() string {
	return fmt.Sprintf("(%d,%d,%d)", pt.x, pt.y, pt.z)
}

func (pt pos) adjacent() []pos {
	return []pos{
		{pt.x - 1, pt.y, pt.z},
		{pt.x + 1, pt.y, pt.z},
		{pt.x, pt.y - 1, pt.z},
		{pt.x, pt.y + 1, pt.z},
		{pt.x, pt.y, pt.z - 1},
		{pt.x, pt.y, pt.z + 1},
	}
}

type puzzle struct {
	data     map[pos]wat
	min, max pos
}

func (p *puzzle) countSurface(target wat) int {
	sum := 0
	for pt := range p.data {
		if p.data[pt] == LAVA {
			for _, pt2 := range pt.adjacent() {
				if p.data[pt2] == target {
					sum++
				}
			}
		}
	}
	return sum
}

func (p *puzzle) fill(pt pos) {
	if pt.x < p.min.x || pt.y < p.min.y || pt.z < p.min.z {
		return
	}
	if pt.x > p.max.x || pt.y > p.max.y || pt.z > p.max.z {
		return
	}

	switch p.data[pt] {
	case AIR:
		p.data[pt] = STEAM
		for _, pt2 := range pt.adjacent() {
			p.fill(pt2)
		}
	case LAVA:
	case STEAM:
	}
}

func (p *puzzle) renderPlane(z int) [][]byte {
	var plane [][]byte
	for y := p.min.y; y <= p.max.y; y++ {
		var line []byte
		for x := p.min.x; x <= p.max.x; x++ {
			pt := pos{x, y, z}
			switch p.data[pt] {
			case AIR:
				line = append(line, '.')
			case LAVA:
				line = append(line, '#')
			case STEAM:
				line = append(line, '*')
			default:
				panic(p.data[pt])
			}
		}
		plane = append(plane, line)
	}
	return plane
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
