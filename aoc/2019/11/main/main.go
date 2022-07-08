package main

import "fmt"

type point struct {
	x, y int
}

func (p point) next(d direction) point {
	switch d {
	case N:
		return point{x: p.x, y: p.y - 1}
	case E:
		return point{x: p.x + 1, y: p.y}
	case S:
		return point{x: p.x, y: p.y + 1}
	case W:
		return point{x: p.x - 1, y: p.y}
	default:
		panic(d)
	}
}

type direction int

const (
	N = direction(iota)
	E
	S
	W
	numDirs
)

func (d direction) String() string {
	return []string{"^", ">", "v", "<"}[d]
}

func main() {
	codes := parseIntCodes(`99`)

	var pos, max point
	pointData := map[point]bool{
		pos: true,
	}
	dir := N

	in := func() int {
		return boolVal(pointData[pos])
	}

	var buf []int
	i := 0
	out := func(v int) {
		buf = append(buf, v)
		if len(buf) < 2 {
			return
		}
		color, face := buf[0], buf[1]
		buf = nil

		pointData[pos] = color == 1

		if face == 0 {
			dir = (dir + numDirs - 1) % numDirs
		} else {
			dir = (dir + 1) % numDirs
		}
		pos = pos.next(dir)
		if max.x < pos.x {
			max.x = pos.x
		}
		if max.y < pos.y {
			max.y = pos.y
		}
		fmt.Printf("%d: %d %s (%d,%d)\n", i, color, dir, pos.x, pos.y)
		i++
	}

	m := NewIntMachine(codes, in, out)
	m.Run()

	fmt.Println(len(pointData))

	for y := 0; y <= max.y; y++ {
		for x := 0; x <= max.x; x++ {
			if pointData[point{x, y}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
