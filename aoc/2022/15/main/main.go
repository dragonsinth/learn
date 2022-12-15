package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var sample = `
Sensor at x=2, y=18: closest beacon is at x=-2, y=15
Sensor at x=9, y=16: closest beacon is at x=10, y=16
Sensor at x=13, y=2: closest beacon is at x=15, y=3
Sensor at x=12, y=14: closest beacon is at x=10, y=16
Sensor at x=10, y=20: closest beacon is at x=10, y=16
Sensor at x=14, y=17: closest beacon is at x=10, y=16
Sensor at x=8, y=7: closest beacon is at x=2, y=10
Sensor at x=2, y=0: closest beacon is at x=2, y=10
Sensor at x=0, y=11: closest beacon is at x=2, y=10
Sensor at x=20, y=14: closest beacon is at x=25, y=17
Sensor at x=17, y=20: closest beacon is at x=21, y=22
Sensor at x=16, y=7: closest beacon is at x=15, y=3
Sensor at x=14, y=3: closest beacon is at x=15, y=3
Sensor at x=20, y=1: closest beacon is at x=15, y=3
`

var (
	re = regexp.MustCompile(`^Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)$`)
)

func main() {
	run(sample, true, 10, pos{20, 20})
	//run(data, false, 2000000, pos{4000000, 4000000})
}

func run(input string, dbg bool, row int, max pos) {
	p := puzzle{
		data: map[pos]byte{},
	}

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !re.MatchString(line) {
			panic(line)
		}

		matches := re.FindStringSubmatch(line)
		beacon := pos{
			x: mustInt(matches[3]),
			y: mustInt(matches[4]),
		}

		p.beacons = append(p.beacons, beacon)
		p.addPoint(beacon, 'B')

		pt := pos{
			x: mustInt(matches[1]),
			y: mustInt(matches[2]),
		}

		dist := manhattan(pt, beacon)
		p.sensors = append(p.sensors, sensor{
			pos:  pt,
			dist: dist,
		})
		p.addPoint(pt, 'S')

		p.ensure(pos{pt.x - dist, pt.y})
		p.ensure(pos{pt.x + dist, pt.y})
		p.ensure(pos{pt.x, pt.y - dist})
		p.ensure(pos{pt.x, pt.y + dist})
	}

	fmt.Println(p.min, p.max)
	if dbg {
		termbox.RenderPlain(p.renderView(), os.Stdout)
		fmt.Println()

		p.fill()

		fmt.Println(p.min, p.max)
		termbox.RenderPlain(p.renderView(), os.Stdout)
		fmt.Println()
	}

	sum := 0
	for x := p.min.x; x <= p.max.x; x++ {
		// check each sensor
		inRange := false
		pt := pos{x, row}
		for _, sn := range p.sensors {
			inRange = inRange || manhattan(pt, sn.pos) <= sn.dist
		}
		if inRange && p.get(pt) != 'B' && p.get(pt) != 'S' {
			sum++
		}
	}
	fmt.Println(sum)

	// Break the entire grid into boxes and traverse from the center of each box
	xPos := []int{0, max.x}
	yPos := []int{0, max.y}
	for _, sn := range p.sensors {
		if sn.x < max.x {
			xPos = append(xPos, sn.x)
		}
		if sn.y < max.y {
			yPos = append(yPos, sn.y)
		}
	}
	xPos = sortUnique(xPos)
	yPos = sortUnique(yPos)
	best, bestDist := pos{}, math.MaxInt
	for i, ic := 0, len(xPos)-1; i < ic; i++ {
		for j, jc := 0, len(yPos)-1; j < jc; j++ {
			min := pos{xPos[i], yPos[j]}
			max := pos{xPos[i+1], yPos[j+1]}
			st := pos{(min.x + max.x) / 2, (min.y + max.y) / 2}
			pt, dist := p.findMinima(dbg, st, min, max)
			if dbg {
				fmt.Println(min, st, max, "=>", pt, dist)
			}
			if dist < bestDist {
				best, bestDist = pt, dist
			}
		}
	}

	if dbg {
		termbox.RenderPlain(p.renderView(), os.Stdout)
		fmt.Println()
	}

	fmt.Println(best, bestDist, "=", 4000000*best.x+best.y)
	fmt.Println(len(xPos)*len(yPos), "regions tested")

	//3172756, 2767556
}

func sortUnique(in []int) []int {
	sort.Ints(in)
	wIdx := 0
	var last int
	for i, v := range in {
		if i == 0 || v != last {
			in[wIdx] = v
			wIdx++
		}
		last = v
	}
	return in[:wIdx]
}

type pos struct {
	x, y int
}

func (pt pos) String() string {
	return fmt.Sprintf("(%d,%d)", pt.x, pt.y)
}

func (pt pos) adjacent() []pos {
	return []pos{
		{pt.x - 1, pt.y},
		{pt.x + 1, pt.y},
		{pt.x, pt.y - 1},
		{pt.x, pt.y + 1},
	}
}

func (pt pos) allAdjacent() []pos {
	return []pos{
		{pt.x - 1, pt.y},
		{pt.x + 1, pt.y},
		{pt.x, pt.y - 1},
		{pt.x, pt.y + 1},
		{pt.x - 1, pt.y - 1},
		{pt.x + 1, pt.y + 1},
		{pt.x + 1, pt.y - 1},
		{pt.x - 1, pt.y + 1},
	}
}

type sensor struct {
	pos
	dist int
}

type puzzle struct {
	sensors []sensor
	beacons []pos

	data map[pos]byte

	min, max pos
}

func (p *puzzle) addPoint(pt pos, wat byte) {
	p.data[pt] = wat
	p.ensure(pt)
}

func (p *puzzle) ensure(pt pos) {
	p.min.x = min(p.min.x, pt.x-1)
	p.max.x = max(p.max.x, pt.x+1)
	p.min.y = min(p.min.y, pt.y-1)
	p.max.y = max(p.max.y, pt.y+1)
}

func (p *puzzle) get(pt pos) byte {
	if c, ok := p.data[pt]; ok {
		return c
	} else {
		return '.'
	}
}

func (p *puzzle) fill() {
	for _, sn := range p.sensors {
		p.doFill(sn)
	}
}

func (p *puzzle) doFill(sn sensor) {
	seen := map[pos]bool{}
	work := []sensor{sn}

	for len(work) > 0 {
		w := work[0]
		work = work[1:]

		pt := w.pos
		if seen[pt] {
			continue
		}
		seen[pt] = true
		if p.get(pt) == '.' {
			p.addPoint(pt, '#')
		}

		if w.dist > 0 {
			for _, pt := range pt.adjacent() {
				work = append(work, sensor{pt, w.dist - 1})
			}
		}
	}
}

func (p *puzzle) renderView() [][]byte {
	var ret [][]byte
	for y := p.min.y; y <= p.max.y; y++ {
		buf := make([]byte, p.max.x-p.min.x+1)
		for i := range buf {
			pt := pos{i + p.min.x, y}
			buf[i] = p.get(pt)
		}
		ret = append(ret, buf)
	}
	return ret
}

func (p *puzzle) findMinima(dbg bool, pt pos, min pos, max pos) (pos, int) {
	seen := map[pos]bool{pt: true}
	bestPt, best := pos{}, math.MaxInt
	work := []pos{pt}
	for len(work) > 0 {
		pt := work[0]
		work = work[1:]

		dist := p.calcDistance(pt)
		if dbg && dist < 10 && p.get(pt) == '#' {
			p.data[pt] = '0' + byte(dist)
		}
		if dist >= best {
			continue // don't even explore this
		}
		if dist == 0 {
			return pt, 0
		}
		if dist < best {
			bestPt, best = pt, dist
			if dbg {
				fmt.Println(bestPt, best)
			}
			seen = map[pos]bool{}
		}
		for _, pt := range pt.allAdjacent() {
			if pt.x >= min.x && pt.x <= max.x && pt.y >= min.y && pt.y <= max.y && !seen[pt] {
				seen[pt] = true
				work = append(work, pt)
			}
		}
	}
	return bestPt, best
}

func (p *puzzle) calcDistance(pt pos) int {
	sum := 0
	for _, sn := range p.sensors {
		dist := manhattan(pt, sn.pos)
		if dist < sn.dist+1 {
			sum += sn.dist + 1 - dist
		}
	}
	return sum
}

func manhattan(a pos, b pos) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
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

func abs(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}
