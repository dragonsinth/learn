package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"math/rand"
	"os"
	"time"
)

type direction int

const (
	N = direction(iota)
	NE
	E
	SE
	S
	SW
	W
	NW
	C
)

func (d direction) isCardinal() bool {
	return d == N || d == S || d == E || d == W
}

type walkOracle struct {
	seen map[pos]int
	sz   int

	coefficients [8]int
	straightVals [8][]int
}

func (wo *walkOracle) estimatePoint(x, y int) int {
	// compute exact near center
	lo := -2 * wo.sz
	hi := 3 * wo.sz
	if x >= lo && x < hi && y >= lo && y < hi {
		return wo.seen[pos{x, y}]
	}

	switch {
	case x < 0 && y < 0:
		return wo.coefficients[NW] + abs(x) + abs(y)
	case x < 0 && y >= wo.sz:
		return wo.coefficients[SW] + abs(x) + abs(y)
	case x >= wo.sz && y < 0:
		return wo.coefficients[NE] + abs(x) + abs(y)
	case x >= wo.sz && y >= wo.sz:
		return wo.coefficients[SE] + abs(x) + abs(y)
	case x < 0:
		return wo.straightVals[W][y] + abs(x) + abs(y)
	case x >= wo.sz:
		return wo.straightVals[E][y] + abs(x) + abs(y)
	case y < 0:
		return wo.straightVals[N][x] + abs(x) + abs(y)
	case y >= wo.sz:
		return wo.straightVals[S][x] + abs(x) + abs(y)
	default:
		panic("here")
	}
}

func getQuadrant(x int, y int) direction {
	switch {
	case x < 0 && y < 0:
		return NW
	case x < 0 && y > 0:
		return SW
	case x > 0 && y < 0:
		return NE
	case x > 0 && y > 0:
		return SE
	case x < 0:
		return W
	case x > 0:
		return E
	case y < 0:
		return N
	case y > 0:
		return S
	default:
		panic("here")
	}
}

func (wo *walkOracle) test() {
	rand.Seed(time.Now().UnixMilli())

	// pick any random x value that's == 0 mod 11, and any random y value
	for i := 0; i < 1000; i++ {
		x := (rand.Intn(10) - 5) * wo.sz
		y := rand.Intn(10*wo.sz) - 5*wo.sz
		if v, ok := wo.seen[pos{x, y}]; ok {
			if v != wo.estimatePoint(x, y) {
				fmt.Println("wrong:", x, y, v, wo.estimatePoint(x, y))
			}
		} else {
			fmt.Println("not seen: ", x, y)
		}
	}

	// pick any random y value that's == 0 mod 11, and any random x value
	for i := 0; i < 1000; i++ {
		y := (rand.Intn(10) - 5) * wo.sz
		x := rand.Intn(10*wo.sz) - 5*wo.sz
		if v, ok := wo.seen[pos{x, y}]; ok {
			if v != wo.estimatePoint(x, y) {
				fmt.Println("wrong:", x, y, v, wo.estimatePoint(x, y))
			}
		} else {
			fmt.Println("not seen: ", x, y)
		}
	}
}

func run2(input string, maxDist int, debug bool) {
	p := parse(input)
	if debug {
		termbox.RenderPlain(p.Render(), os.Stdout)
		fmt.Println()
	}

	if maxDist < 20*p.sz {
		fmt.Print("real: ")
		seen := p.walk(maxDist)
		fmt.Println(countWalkable(seen, maxDist))
		return
	}

	st := time.Now()
	defer func() {
		fmt.Println(time.Since(st))
	}()

	seen := p.walk(20 * p.sz)
	wo := &walkOracle{
		seen: seen,
		sz:   p.sz,
	}

	d, d2 := p.sz*5, p.sz*10
	wo.coefficients[NW] = seen[pos{-d, -d}] - d2
	wo.coefficients[NE] = seen[pos{d, -d}] - d2
	wo.coefficients[SE] = seen[pos{d, d}] - d2
	wo.coefficients[SW] = seen[pos{d, -d}] - d2

	calcXVals := func(y int) []int {
		var ret []int
		for x := 0; x < p.sz; x++ {
			v := seen[pos{x, y}]
			ret = append(ret, v-abs(x)-abs(y))
		}
		return ret
	}
	calcYVals := func(x int) []int {
		var ret []int
		for y := 0; y < p.sz; y++ {
			v := seen[pos{x, y}]
			ret = append(ret, v-abs(x)-abs(y))
		}
		return ret
	}

	wo.straightVals[N] = calcXVals(-d)
	wo.straightVals[S] = calcXVals(d)
	wo.straightVals[E] = calcYVals(d)
	wo.straightVals[W] = calcYVals(-d)

	wo.test()
	fmt.Println("done testing")

	findBound := func(step int, swapXY bool) int {
		if swapXY {
			i := 0
			for {
				y := p.sz * i
				if wo.estimatePoint(0, y) > maxDist && wo.estimatePoint(p.sz-1, y) > maxDist {
					return i
				}
				i += step
			}
		} else {
			i := 0
			for {
				x := p.sz * i
				if wo.estimatePoint(x, 0) > maxDist && wo.estimatePoint(x, p.sz-1) > maxDist {
					return i
				}
				i += step
			}
		}
	}

	txmin := findBound(-1, false)
	txmax := findBound(1, false)
	tymin := findBound(-1, true)
	tymax := findBound(1, true)
	tDistMax := max(max(max(-txmin, txmax), -tymin), tymax)
	fmt.Println(tDistMax)

	evenFull := p.walkTile(wo, 2, 2, maxDist)
	oddFull := p.walkTile(wo, 2, 3, maxDist)
	fmt.Println(evenFull, oddFull)

	xpy := map[int]int{}
	xny := map[int]int{}

	// iterate every possible tile adding its score
	sum := 0
	for ty := -tDistMax; ty < tDistMax; ty++ {
		for tx := -tDistMax; tx < tDistMax; tx++ {
			tdist := abs(tx) + abs(ty)
			if tdist > tDistMax {
				continue // not reachable
			}
			if tDistMax-tdist > 2 {
				// internal tile
				if tdist%2 == 0 {
					sum += evenFull
				} else {
					sum += oddFull
				}
				continue
			}

			// check the cache
			quad := getQuadrant(tx, ty)
			switch quad {
			case NW, SE:
				v, ok := xpy[tx+ty]
				if !ok {
					v = p.walkTile(wo, tx, ty, maxDist)
					xpy[tx+ty] = v
				}
				sum += v
			case SW, NE:
				v, ok := xny[tx-ty]
				if !ok {
					v = p.walkTile(wo, tx, ty, maxDist)
					xny[tx-ty] = v
				}
				sum += v
			case N, S, E, W:
				sum += p.walkTile(wo, tx, ty, maxDist)
			default:
				panic("here")
			}
		}
	}

	fmt.Println(sum)
}

func (p puz) walkTile(wo *walkOracle, tx int, ty int, maxDist int) int {
	minX, minY := tx*p.sz, ty*p.sz
	maxX, maxY := minX+p.sz-1, minY+p.sz-1
	type task struct {
		pos
		dist int
	}
	var work []task
	seen := map[pos]int{}

	addTask := func(x, y, dist int) {
		pt := pos{x, y}
		seen[pt] = dist
		work = append(work, task{pt, dist})
	}

	quad := getQuadrant(tx, ty)
	switch quad {
	case N:
		// add the entire south edge
		y := maxY
		for x := minX; x <= maxX; x++ {
			addTask(x, y, wo.estimatePoint(x, y))
		}
	case NE:
		// add just the SW corner
		x, y := minX, maxY
		addTask(x, y, wo.estimatePoint(x, y))
	case E:
		// add the entire west edge
		x := minX
		for y := minY; y <= maxY; y++ {
			addTask(x, y, wo.estimatePoint(x, y))
		}
	case SE:
		// add just the NW corner
		x, y := minX, minY
		addTask(x, y, wo.estimatePoint(x, y))
	case S:
		// add the entire north edge
		y := minY
		for x := minX; x <= maxX; x++ {
			addTask(x, y, wo.estimatePoint(x, y))
		}
	case SW:
		// add just the NE corner
		x, y := maxX, minY
		addTask(x, y, wo.estimatePoint(x, y))
	case W:
		// add the entire east edge
		x := maxX
		for y := minY; y <= maxY; y++ {
			addTask(x, y, wo.estimatePoint(x, y))
		}
	case NW:
		// add just the SE corner
		x, y := maxX, maxY
		addTask(x, y, wo.estimatePoint(x, y))
	default:
		panic("here")
	}

	maybeWalk := func(x, y, dist int) {
		if x < minX || x > maxX || y < minY || y > maxY {
			return // out of bounds
		}

		pt := pos{x, y}
		if v, ok := seen[pt]; ok {
			if v <= dist {
				return
			}
		}
		if p.at(x, y) == '.' {
			addTask(x, y, dist)
		}
	}

	for len(work) > 0 {
		t := work[0]
		work = work[1:]
		x, y, dist := t.x, t.y, t.dist+1
		maybeWalk(x, y-1, dist)
		maybeWalk(x+1, y, dist)
		maybeWalk(x, y+1, dist)
		maybeWalk(x-1, y, dist)
	}

	sum := 0
	for _, v := range seen {
		if v <= maxDist && v%2 == maxDist%2 {
			sum++
		}
	}
	return sum
}
