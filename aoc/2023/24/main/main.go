package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var sample = `
19, 13, 30 @ -2,  1, -2
18, 19, 22 @ -1, -1, -2
20, 25, 34 @ -2, -2, -4
12, 31, 28 @ -1, -2, -1
20, 19, 15 @  1, -5, -3
`

var re = regexp.MustCompile(`^(\d+),\s*(\d+),\s*(\d+)\s*@\s*(-?\d+),\s*(-?\d+),\s*(-?\d+)$`)

func main() {
	part1(sample, 7, 27, true)
	//part1(data, 200000000000000, 400000000000000, false)
}

func part1(input string, min int, max int, debug bool) {
	p := parse(input)
	mn, mx := float64(min), float64(max)
	sum := 0
	for i := range p.vecs {
		for j := range p.vecs {
			if i < j {
				if debug {
					fmt.Println(i, p.vecs[i].String())
					fmt.Println(j, p.vecs[j].String())
				}
				x1, y1, ok1 := intersectFloat2d(p.vecs[i], p.vecs[j], X, Y, debug)
				x2, y2, ok2 := intersectFloat2d(p.vecs[j], p.vecs[i], X, Y, debug)
				if debug {
					fmt.Println(x1, y1, ok1)
					fmt.Println(x2, y2, ok2)
				}
				if !ok1 || !ok2 {
					continue
				}
				if !eq(x1, x2) || !eq(y1, y2) {
					panic("here")
				}
				if x1 < mn || x1 > mx {
					continue
				}
				if y1 < mn || y1 > mx {
					continue
				}
				sum++
			}
		}
	}
	fmt.Println(sum)
}

func parse(input string) puz {
	var p puz
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !re.MatchString(line) {
			panic(line)
		}
		matches := re.FindStringSubmatch(line)
		p.vecs = append(p.vecs, vec{
			p: [3]int{
				mustInt(matches[1]),
				mustInt(matches[2]),
				mustInt(matches[3]),
			},
			v: [3]int{
				mustInt(matches[4]),
				mustInt(matches[5]),
				mustInt(matches[6]),
			},
		})
	}
	return p
}

type puz struct {
	vecs []vec
}

type dim int

func (d dim) String() string {
	return "XYZ"[d : d+1]
}

const (
	X = dim(iota)
	Y
	Z
)

type vec struct {
	p [3]int
	v [3]int
}

func (v vec) String() any {
	return fmt.Sprintf("%d, %d, %d @ %d, %d, %d", v.p[X], v.p[Y], v.p[Z], v.v[X], v.v[Y], v.v[Z])
}

func intersectFloat2d(v1 vec, v2 vec, d1 dim, d2 dim, debug bool) (float64, float64, bool) {
	a, b, x, y := float64(v1.p[d1]), float64(v1.p[d2]), float64(v1.v[d1]), float64(v1.v[d2])
	c, d, t, u := float64(v2.p[d1]), float64(v2.p[d2]), float64(v2.v[d1]), float64(v2.v[d2])

	dt := (u*c - u*a + t*b - t*d) / (u*x - t*y)
	if math.IsNaN(dt) {
		if debug {
			fmt.Println("no intersection")
		}
		return 0, 0, false
	}
	if dt < 0 {
		if debug {
			fmt.Println("past intersection")
		}
		return 0, 0, false
	}
	return a + x*dt, b + y*dt, true
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

func eq(a float64, b float64) bool {
	const Epsilon = 1e-9
	diff := math.Abs(a - b)
	if diff <= Epsilon*math.Abs(a) || diff <= Epsilon*math.Abs(b) {
		return true
	}
	return false
}
