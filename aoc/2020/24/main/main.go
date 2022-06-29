package main

import (
	"fmt"
	"regexp"
	"strings"
)

var sample = `
sesenwnenenewseeswwswswwnenewsewsw
neeenesenwnwwswnenewnwwsewnenwseswesw
seswneswswsenwwnwse
nwnwneseeswswnenewneswwnewseswneseene
swweswneswnenwsewnwneneseenw
eesenwseswswnenwswnwnwsewwnwsene
sewnenenenesenwsewnenwwwse
wenwwweseeeweswwwnwwe
wsweesenenewnwwnwsenewsenwwsesesenwne
neeswseenwwswnwswswnw
nenwswwsewswnenenewsenwsenwnesesenew
enewnwewneswsewnwswenweswnenwsenwsw
sweneswneswneneenwnewenewwneswswnese
swwesenesewenwneswnwwneseswwne
enesenwswwswneneswsenwnewswseenwsese
wnwnesenesenenwwnenwsewesewsesesew
nenewswnwewswnenesenwnesewesw
eneswnwswnwsenenwnwnwwseeswneewsenese
neswnwewnwnwseenwseesewsenwsweewe
wseweeenwnesenwwwswnew
`

var input = sample

var parseLine = regexp.MustCompile(`^[nsew]+$`)

type point struct {
	q, r, s int
}

func (p point) neighbors() []point {
	return []point{
		p.move(00, +1, -1),
		p.move(00, -1, +1),
		p.move(+1, 00, -1),
		p.move(-1, 00, +1),
		p.move(+1, -1, 00),
		p.move(-1, +1, 00),
	}
}

func (p point) move(q int, r int, s int) point {
	return point{
		q: p.q + q,
		r: p.r + r,
		s: p.s + s,
	}
}

func main() {
	lit := map[point]bool{}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !parseLine.MatchString(line) {
			panic(line)
		}

		var p point
		for line != "" {
			if c := line[0]; c == 'e' || c == 'w' {
				line = line[1:]
				switch c {
				case 'e':
					p.q++
					p.s--
				case 'w':
					p.q--
					p.s++
				}
			} else {
				pre := line[0:2]
				line = line[2:]
				switch pre {
				case "ne":
					p.q++
					p.r--
				case "sw":
					p.q--
					p.r++
				case "nw":
					p.r--
					p.s++
				case "se":
					p.r++
					p.s--
				default:
					panic("should not")
				}
			}
		}

		if lit[p] {
			delete(lit, p)
		} else {
			lit[p] = true
		}
	}

	fmt.Println(len(lit))

	for i := 0; i < 100; i++ {
		lit = nextStep(lit)
		fmt.Printf("Day %d: %d\n", i+1, len(lit))
	}
}

func nextStep(lit map[point]bool) map[point]bool {
	ret := map[point]bool{}

	// For each point in the input, compute both it and all surrounding points.
	for p := range lit {
		compute(p, lit, ret)
		for _, pn := range p.neighbors() {
			compute(pn, lit, ret)
		}
	}

	// Filter
	for k := range ret {
		if !ret[k] {
			delete(ret, k)
		}
	}
	return ret
}

func compute(p point, lit map[point]bool, ret map[point]bool) {
	_, ok := ret[p]
	if ok {
		return // already computed
	}

	count := 0
	for _, pn := range p.neighbors() {
		if lit[pn] {
			count++
		}
	}

	if lit[p] {
		ret[p] = count == 1 || count == 2
	} else {
		ret[p] = count == 2
	}
}
