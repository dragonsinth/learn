package main

import "fmt"

type puzzle struct {
	v      [][][]byte
	width  int
	height int
	depth  int
}

func (p *puzzle) Active() int {
	sum := 0
	for _, plane := range p.v {
		for _, row := range plane {
			for _, col := range row {
				if col != 0 {
					sum++
				}
			}
		}
	}
	return sum
}

func (p *puzzle) Print() {
	for z, plane := range p.v {
		fmt.Printf("z=%d\n", z-p.depth/2)
		for _, row := range plane {
			for _, col := range row {
				fmt.Print(string(".#"[col]))
			}
			fmt.Println()
		}
	}
	fmt.Println(p.Active())
}

func (p *puzzle) Next() *puzzle {
	next := puzzle{
		v:      make([][][]byte, p.depth+2),
		width:  p.width + 2,
		height: p.height + 2,
		depth:  p.depth + 2,
	}

	// First create all the output objects.
	for i := range next.v {
		next.v[i] = make([][]byte, next.height)
		for j := range next.v[i] {
			next.v[i][j] = make([]byte, next.width)
		}
	}

	// Now compute the output from the input
	for z := 0; z < next.depth; z++ {
		for y := 0; y < next.height; y++ {
			for x := 0; x < next.width; x++ {
				sum := p.sum(x-1, y-1, z-1)
				active := p.active(x-1, y-1, z-1)
				if active && (sum == 2 || sum == 3) {
					next.v[z][y][x] = 1
				} else if !active && sum == 3 {
					next.v[z][y][x] = 1
				}
			}
		}
	}

	return &next
}

func (p *puzzle) active(x int, y int, z int) bool {
	if x < 0 || y < 0 || z < 0 {
		return false
	}
	if x >= p.width || y >= p.height || z >= p.depth {
		return false
	}
	return p.v[z][y][x] != 0
}

func (p *puzzle) sum(x int, y int, z int) int {
	sum := 0
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			for k := z - 1; k <= z+1; k++ {
				if i == x && j == y && k == z {
					continue
				}
				if p.active(i, j, k) {
					sum++
				}
			}
		}
	}
	return sum
}
