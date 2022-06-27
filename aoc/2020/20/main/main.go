package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	size  = 10
	isize = 8
)

var (
	input = sample

	headerRe = regexp.MustCompile(`^Tile (\d+):$`)
	bodyRe   = regexp.MustCompile(`^[.#]{10}$`)
)

// Denotes the orientation of the original north edge.
type orient int

func (o orient) String() string {
	switch o {
	case NORT:
		return "NORT"
	case EAST:
		return "EAST"
	case SOUT:
		return "SOUT"
	case WEST:
		return "WEST"
	case TRON:
		return "TRON"
	case TSAE:
		return "TSAE"
	case TUOS:
		return "TUOS"
	case TSEW:
		return "TSEW"
	default:
		panic(o)
	}
}

const (
	NORT = orient(iota)
	EAST
	SOUT
	WEST
	TRON
	TSAE
	TUOS
	TSEW
	MAX_ORIENT
)

type puzzle struct {
	size  int
	tiles []*tile
	byId  map[int]*tile
}

type tileType int

func (t tileType) String() string {
	switch t {
	case CENTER:
		return "CENTER"
	case EDGE:
		return "EDGE"
	case CORNER:
		return "CORNER"
	default:
		panic(t)
	}
}

const (
	CENTER = tileType(iota)
	EDGE
	CORNER
)

type tile struct {
	id   int
	bits [][]bool

	edges [MAX_ORIENT][4]uint16
	typ   tileType
}

func (t *tile) InitEdges() {
	// All bits are read clockwise; mirror images read counterclockwise.
	n, np := t.rowToBits(0)
	e, ep := t.colToBits(9)
	sp, s := t.rowToBits(9)
	wp, w := t.colToBits(0)

	t.edges[NORT] = [4]uint16{n, e, s, w}
	t.edges[EAST] = [4]uint16{w, n, e, s}
	t.edges[SOUT] = [4]uint16{s, w, n, e}
	t.edges[WEST] = [4]uint16{e, s, w, n}

	t.edges[TRON] = [4]uint16{np, wp, sp, ep}
	t.edges[TSAE] = [4]uint16{ep, np, wp, sp}
	t.edges[TUOS] = [4]uint16{sp, ep, np, wp}
	t.edges[TSEW] = [4]uint16{wp, sp, ep, np}
}

func (t *tile) rowToBits(r int) (uint16, uint16) {
	var l2r, r2l uint16
	for i := 0; i < size; i++ {
		l2r <<= 1
		if t.bits[r][i] {
			l2r |= 1
		}
		r2l <<= 1
		if t.bits[r][size-i-1] {
			r2l |= 1
		}

	}
	return l2r, r2l
}

func (t *tile) colToBits(c int) (uint16, uint16) {
	var t2b, b2t uint16
	for i := 0; i < size; i++ {
		t2b <<= 1
		if t.bits[i][c] {
			t2b |= 1
		}
		b2t <<= 1
		if t.bits[size-i-1][c] {
			b2t |= 1
		}

	}
	return t2b, b2t
}

func (t *tile) Render(or orient) [][]bool {
	ret := makeBoolField(isize, isize)
	rot := rotate(t.bits, or)
	for y := 0; y < isize; y++ {
		for x := 0; x < isize; x++ {
			ret[y][x] = rot[y+1][x+1]
		}
	}
	return ret
}

func main() {
	puz := puzzle{
		byId: map[int]*tile{},
	}

	var curTile *tile
	curLine := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if headerRe.MatchString(line) {
			if curTile != nil {
				panic("should be nil")
			}
			curTile = &tile{
				bits: makeBoolField(size, size),
			}
			curLine = 0
			matches := headerRe.FindStringSubmatch(line)
			curTile.id = mustInt(matches[1])
		} else if bodyRe.MatchString(line) {
			for i, c := range line {
				curTile.bits[curLine][i] = c == '#'
			}
			curLine++
			if curLine == 10 {
				puz.tiles = append(puz.tiles, curTile)
				puz.byId[curTile.id] = curTile
				curTile = nil
				curLine = 0
			}
		} else {
			panic(line)
		}
	}

	if curTile != nil {
		panic("should be nil")
	}

	edgeCounts := map[uint16]int{}
	for _, t := range puz.tiles {
		t.InitEdges()
		for _, e := range t.edges[NORT] {
			edgeCounts[e]++
		}
		for _, e := range t.edges[TRON] {
			edgeCounts[e]++
		}
	}

	// Classify each tile
	for _, t := range puz.tiles {
		uniq := 0
		for _, e := range t.edges[NORT] {
			switch edgeCounts[e] {
			case 1:
				uniq++
			case 2:
			default:
				panic(edgeCounts[e])
			}
		}

		switch uniq {
		case 0:
			t.typ = CENTER
		case 1:
			t.typ = EDGE
		case 2:
			t.typ = CORNER
		}
	}

	switch len(puz.tiles) {
	case 9:
		puz.size = 3
	case 144:
		puz.size = 12
	default:
		panic(len(puz.tiles))
	}

	sol := puz.Solve(edgeCounts)
	fmt.Println(sol)
	fmt.Println(sol.Solve())
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
