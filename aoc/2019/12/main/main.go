package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var samples = []string{
	`
<x=-1, y=0, z=2>
<x=2, y=-10, z=-7>
<x=4, y=-8, z=8>
<x=3, y=5, z=-1>`,
	`
<x=-8, y=-10, z=0>
<x=5, y=5, z=10>
<x=2, y=-7, z=3>
<x=9, y=-8, z=-3>`,
}

var (
	re = regexp.MustCompile(`^<x=(-?\d+), y=(-?\d+), z=(-?\d+)>$`)
)

type moon struct {
	x, y, z    int
	dx, dy, dz int
}

func (m moon) String() string {
	return fmt.Sprintf(`pos=<x=%3d, y=%3d, z=%3d>, vel=<x=%3d, y=%3d, z=%3d>`, m.x, m.y, m.z, m.dx, m.dy, m.dz)
}

func main() {
	newPuzzle(samples[0]).run(10)
	newPuzzle(samples[1]).run(100)

	for _, input := range samples {
		p := newPuzzle(input)
		lx := p.sliceX().findLoop()
		ly := p.sliceY().findLoop()
		lz := p.sliceZ().findLoop()
		fmt.Println(lx, ly, lz, lcm(lx, lcm(ly, lz)))
	}
}

func newPuzzle(input string) *puzzle {
	var p puzzle

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !re.MatchString(line) {
			panic(line)
		}

		parts := re.FindStringSubmatch(line)
		if len(parts) != 4 {
			panic(parts)
		}

		p.moons = append(p.moons, &moon{
			x: mustInt(parts[1]),
			y: mustInt(parts[2]),
			z: mustInt(parts[3]),
		})
	}

	return &p
}

type puzzle struct {
	moons []*moon
}

func (p *puzzle) run(steps int) {
	for step := 0; true; step++ {
		if step%(steps/10) == 0 {
			fmt.Printf("After %d steps:\n", step)
			for _, m := range p.moons {
				fmt.Println(m)
			}
			p, k, t := p.energy()
			fmt.Printf("pot: %d; kin: %d; tot: %d\n", p, k, t)
			fmt.Println()
		}
		if step == steps {
			break
		}

		p.applyGravity()
		p.applyVelocity()
	}

	// compute energy
	sum := 0
	for _, m := range p.moons {
		pot := abs(m.x) + abs(m.y) + abs(m.z)
		kin := abs(m.dx) + abs(m.dy) + abs(m.dz)
		tot := pot * kin
		fmt.Printf("pot: %d; kin: %d; tot: %d\n", pot, kin, tot)
		sum += tot
	}
	fmt.Println("total energy:", sum)
	fmt.Println()
}

func (p *puzzle) applyGravity() {
	for i, m1 := range p.moons {
		for j, m2 := range p.moons {
			// Unique pairs only
			if i < j {
				if m1.x < m2.x {
					m1.dx++
					m2.dx--
				} else if m1.x > m2.x {
					m1.dx--
					m2.dx++
				}
				if m1.y < m2.y {
					m1.dy++
					m2.dy--
				} else if m1.y > m2.y {
					m1.dy--
					m2.dy++
				}
				if m1.z < m2.z {
					m1.dz++
					m2.dz--
				} else if m1.z > m2.z {
					m1.dz--
					m2.dz++
				}
			}
		}
	}
}

func (p *puzzle) applyVelocity() {
	for _, m := range p.moons {
		m.x += m.dx
		m.y += m.dy
		m.z += m.dz
	}
}

func (p *puzzle) energy() (int, int, int) {
	sPot, sKin, sTot := 0, 0, 0
	for _, m := range p.moons {
		pot := abs(m.x) + abs(m.y) + abs(m.z)
		kin := abs(m.dx) + abs(m.dy) + abs(m.dz)
		tot := pot * kin
		sPot += pot
		sKin += kin
		sTot += tot
	}
	return sPot, sKin, sTot
}

func (p *puzzle) sliceX() *puzSlice {
	var ret puzSlice
	for i := 0; i < 4; i++ {
		ret.p[i] = p.moons[i].x
	}
	return &ret
}

func (p *puzzle) sliceY() *puzSlice {
	var ret puzSlice
	for i := 0; i < 4; i++ {
		ret.p[i] = p.moons[i].y
	}
	return &ret
}

func (p *puzzle) sliceZ() *puzSlice {
	var ret puzSlice
	for i := 0; i < 4; i++ {
		ret.p[i] = p.moons[i].z
	}
	return &ret
}

type puzSlice struct {
	p [4]int
	v [4]int
}

func (p puzSlice) findLoop() int {
	seen := map[puzSlice]bool{}
	i := 0
	for {
		if seen[p] {
			return i
		}
		seen[p] = true
		p = p.next()
		i++
	}
}

func (p puzSlice) next() puzSlice {
	// apply gravity
	for i := range p.p {
		for j := range p.p {
			// Unique pairs only
			if i < j {
				if p.p[i] < p.p[j] {
					p.v[i]++
					p.v[j]--
				} else if p.p[i] > p.p[j] {
					p.v[i]--
					p.v[j]++
				}
			}
		}
	}
	// apply velocity
	for i := range p.p {
		p.p[i] += p.v[i]
	}
	return p
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func lcm(a, b int) int {
	return a * (b / gcd(a, b))
}
