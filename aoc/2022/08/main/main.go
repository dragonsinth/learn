package main

import (
	"fmt"
	"strings"
)

var sample = `
30373
25512
65332
33549
35390
`

func main() {
	var p puzzle
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if p.xmax == 0 {
			p.xmax = len(line)
		} else if p.xmax != len(line) {
			panic(line)
		}

		var ln []int
		for _, c := range line {
			ln = append(ln, int(c-'0'))
		}
		p.grid = append(p.grid, ln)
	}
	p.ymax = len(p.grid)
	fmt.Println(p.visible())
	fmt.Println(p.bestScore())
}

type puzzle struct {
	grid       [][]int
	xmax, ymax int
}

func (p puzzle) bestScore() int {
	var best, bestX, bestY int
	for y := 0; y < p.ymax; y++ {
		for x := 0; x < p.xmax; x++ {
			score := p.score(x, y)
			if score > best {
				best, bestX, bestY = score, x, y
			}
		}
	}
	fmt.Println(best, "at", bestX, bestY)
	return best
}

func (p puzzle) score(px, py int) int {
	height := p.grid[py][px]

	// look left
	left := 0
	for x := px - 1; x >= 0; x-- {
		left++
		if p.grid[py][x] >= height {
			break
		}
	}

	// look right
	right := 0
	for x := px + 1; x < p.xmax; x++ {
		right++
		if p.grid[py][x] >= height {
			break
		}
	}

	// look up
	up := 0
	for y := py - 1; y >= 0; y-- {
		up++
		if p.grid[y][px] >= height {
			break
		}
	}

	// look down
	down := 0
	for y := py + 1; y < p.ymax; y++ {
		down++
		if p.grid[y][px] >= height {
			break
		}
	}

	ret := left * right * up * down
	fmt.Println(left, right, up, down, "=", ret)
	return ret
}

func (p puzzle) visible() int {
	var vis [][]bool
	for _, row := range p.grid {
		vis = append(vis, make([]bool, len(row)))
	}

	// left
	for y := 0; y < p.ymax; y++ {
		max := -1
		for x := 0; x < p.xmax; x++ {
			v := p.grid[y][x]
			if v > max {
				vis[y][x] = true
				max = v
			}
		}
	}

	// right
	for y := 0; y < p.ymax; y++ {
		max := -1
		for x := p.xmax - 1; x >= 0; x-- {
			v := p.grid[y][x]
			if v > max {
				vis[y][x] = true
				max = v
			}
		}
	}

	// top
	for x := 0; x < p.xmax; x++ {
		max := -1
		for y := 0; y < p.ymax; y++ {
			v := p.grid[y][x]
			if v > max {
				vis[y][x] = true
				max = v
			}
		}
	}

	// bottom
	for x := 0; x < p.xmax; x++ {
		max := -1
		for y := p.ymax - 1; y >= 0; y-- {
			v := p.grid[y][x]
			if v > max {
				vis[y][x] = true
				max = v
			}
		}
	}

	sum := 0
	for x := 0; x < p.xmax; x++ {
		for y := 0; y < p.ymax; y++ {
			if vis[y][x] {
				sum++
			}
		}
	}

	return sum
}
