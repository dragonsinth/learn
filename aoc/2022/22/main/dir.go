package main

type pos struct {
	x, y int
}

func (p pos) advance(d dir) pos {
	switch d {
	case EAST:
		p.x++
	case SOUTH:
		p.y++
	case WEST:
		p.x--
	case NORTH:
		p.y--
	default:
		panic(d)
	}
	return p
}

type flatPlayer struct {
	pos
	d dir
}

type dir int

const (
	EAST = dir(iota)
	SOUTH
	WEST
	NORTH
)

func (d dir) turn(t turn) dir {
	switch t {
	case 'L':
		return (4 + d - 1) % 4
	case 'R':
		return (d + 1) % 4
	default:
		panic(t)
	}
}

func (d dir) byte() byte {
	return ">v<^"[d]
}

func (d dir) score() int {
	return int(d)
}

type turn byte

func (t turn) String() string {
	return string(t)
}
