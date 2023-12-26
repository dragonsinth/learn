package main

import (
	"fmt"
	"math"
)

func part2(input string, min int, max int, debug bool) {
	p := parse(input)
	var sum [3]int
	for dx := min; dx <= max; dx++ {
		for dy := min; dy <= max; dy++ {
			px, py, ok := p.allIntersect(float64(dx), float64(dy), X, Y, false)
			if ok {
				fmt.Println(dx, dy, px, py, ok)
				sum[X] = px
				sum[Y] = py
			}
		}
	}
	for dx := min; dx <= max; dx++ {
		for dy := min; dy <= max; dy++ {
			px, py, ok := p.allIntersect(float64(dx), float64(dy), X, Z, false)
			if ok {
				fmt.Println(dx, dy, px, py, ok)
				sum[X] = px
				sum[Z] = py
			}
		}
	}
	for dx := min; dx <= max; dx++ {
		for dy := min; dy <= max; dy++ {
			px, py, ok := p.allIntersect(float64(dx), float64(dy), Y, Z, false)
			if ok {
				fmt.Println(dx, dy, px, py, ok)
				sum[Y] = px
				sum[Z] = py
			}
		}
	}
	fmt.Println(sum[X] + sum[Y] + sum[Z])
}

func (p puz) allIntersect(dx, dy float64, d1 dim, d2 dim, debug bool) (int, int, bool) {
	hit, miss := 0, 0
	seen := false
	rx, ry := 0.0, 0.0
	for i, v1 := range p.vecs {
		a := v1
		a.v[d1] -= dx
		a.v[d2] -= dy
		if a.v[d1] == 0 && a.v[d2] == 0 {
			continue // unmoving
		}
		for j, v2 := range p.vecs {
			if i >= j {
				continue
			}
			b := v2
			b.v[d1] -= dx
			b.v[d2] -= dy
			if b.v[d1] == 0 && b.v[d2] == 0 {
				continue // unmoving
			}
			px, py, ok := intersectFloat2d(a, b, d1, d2, debug)
			if debug {
				fmt.Println(a, b, px, py, ok)
			}
			if !ok {
				miss++
				if miss > 3 {
					return 0, 0, false
				}
				continue
			}
			if math.IsNaN(px) {
				hit++
				continue // convergence
			}
			if !seen {
				seen = true
				rx, ry = px, py
				hit++
			} else {
				if !eq(rx, px) || !eq(ry, py) {
					miss++
					if miss > 3 {
						return 0, 0, false
					}
				} else {
					hit++
				}
			}
		}
	}
	if hit > miss*10 {
		if miss > 0 {
			fmt.Println(hit, miss)
		}
	}
	return int(math.Round(rx)), int(math.Round(ry)), hit > miss*10
}
