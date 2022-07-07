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
	codes := parseIntCodes("99")
	m := NewIntMachine(codes)
	go m.Run()

	var pos, max point
	pointData := map[point]bool{
		pos: true,
	}
	dir := N

	for i := 0; true; i++ {
		select {
		case m.Writer() <- boolVal(pointData[pos]):
		case _, ok := <-m.Reader():
			if !ok {
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
				return // done
			}
			panic("should not get here")
		}

		color, face := <-m.Reader(), <-m.Reader()
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
	}
}
