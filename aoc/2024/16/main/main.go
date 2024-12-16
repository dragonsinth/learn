package main

const sample1 = `
###############
#.......#....E#
#.#.###.#.###.#
#.....#.#...#.#
#.###.#####.#.#
#.#.#.......#.#
#.#.#####.###.#
#...........#.#
###.#.#####.#.#
#...#.....#.#.#
#.#.#.###.#.#.#
#.....#...#.#.#
#.###.#.#.#.#.#
#S..#.....#...#
###############
`

const sample2 = `
#################
#...#...#...#..E#
#.#.#.#.#.#.#.#.#
#.#.#.#...#...#.#
#.#.#.#.###.#.#.#
#...#.#.#.....#.#
#.#.#.#.#.#####.#
#.#...#.#.#.....#
#.#.#####.#.###.#
#.#.#.......#...#
#.#.###.#####.###
#.#.#...#.....#.#
#.#.#.#####.###.#
#.#.#.........#.#
#.#.#.#########.#
#S#.............#
#################
`

func main() {
	run(sample1)
	run(sample2)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type pos struct {
	x, y int
}

func (p pos) next(d dir) pos {
	switch d {
	case N:
		return pos{p.x + 0, p.y - 1}
	case E:
		return pos{p.x + 1, p.y + 0}
	case S:
		return pos{p.x + 0, p.y + 1}
	case W:
		return pos{p.x - 1, p.y + 0}
	default:
		panic(d)
	}
}

type dir byte

func (d dir) left() dir {
	return (d + 3) % 4
}

func (d dir) right() dir {
	return (d + 1) % 4
}

const (
	N = dir(iota)
	E
	S
	W
)
