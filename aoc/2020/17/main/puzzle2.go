package main

import "fmt"

type puzzle2 struct {
	v      [][][][]byte
	width  int
	height int
	depth  int
	zonk   int
}

func (p *puzzle2) Active() int {
	sum := 0
	for _, cube := range p.v {
		for _, plane := range cube {
			for _, row := range plane {
				for _, col := range row {
					if col != 0 {
						sum++
					}
				}
			}
		}
	}
	return sum
}

func (p *puzzle2) Print() {
	for w, cube := range p.v {
		for z, plane := range cube {
			fmt.Printf("z=%d, w=%d\n", z-p.depth/2, w-p.zonk/2)
			for _, row := range plane {
				for _, col := range row {
					fmt.Print(string(".#"[col]))
				}
				fmt.Println()
			}
		}
	}
	fmt.Println(p.Active())
}

func (p *puzzle2) Next() *puzzle2 {
	next := puzzle2{
		v:      make([][][][]byte, p.zonk+2),
		width:  p.width + 2,
		height: p.height + 2,
		depth:  p.depth + 2,
		zonk:   p.zonk + 2,
	}

	// First create all the output objects.
	for i := range next.v {
		next.v[i] = make([][][]byte, next.depth)
		for j := range next.v[i] {
			next.v[i][j] = make([][]byte, next.height)
			for k := range next.v[i][j] {
				next.v[i][j][k] = make([]byte, next.width)
			}
		}
	}

	// Now compute the output from the input
	for w := 0; w < next.zonk; w++ {
		for z := 0; z < next.depth; z++ {
			for y := 0; y < next.height; y++ {
				for x := 0; x < next.width; x++ {
					sum := p.sum(x-1, y-1, z-1, w-1)
					active := p.active(x-1, y-1, z-1, w-1)
					if active && (sum == 2 || sum == 3) {
						next.v[w][z][y][x] = 1
					} else if !active && sum == 3 {
						next.v[w][z][y][x] = 1
					}
				}
			}
		}
	}

	return &next
}

func (p *puzzle2) active(x int, y int, z int, w int) bool {
	if x < 0 || y < 0 || z < 0 || w < 0 {
		return false
	}
	if x >= p.width || y >= p.height || z >= p.depth || w >= p.zonk {
		return false
	}
	return p.v[w][z][y][x] != 0
}

func (p *puzzle2) sum(x int, y int, z int, w int) int {
	sum := 0
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			for k := z - 1; k <= z+1; k++ {
				for l := w - 1; l <= w+1; l++ {
					if i == x && j == y && k == z && l == w {
						continue
					}
					if p.active(i, j, k, l) {
						sum++
					}
				}
			}
		}
	}
	return sum
}
