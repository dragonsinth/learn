package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `
R 4
U 4
L 3
D 1
R 4
D 1
L 5
R 2
`

var sample2 = `
R 5
U 8
L 8
D 3
R 17
D 10
L 25
U 20
`

type pos struct {
	x, y int
}

func (p pos) move(dx int, dy int) pos {
	return pos{x: p.x + dx, y: p.y + dy}
}

func (p pos) moveTowards(head pos) pos {
	dx := head.x - p.x
	dy := head.y - p.y
	if abs(dx) <= 1 && abs(dy) <= 1 {
		return p // in range
	}

	ret := p
	if dx < 0 {
		ret.x--
	} else if dx > 0 {
		ret.x++
	}

	if dy < 0 {
		ret.y--
	} else if dy > 0 {
		ret.y++
	}

	return ret
}

func main() {
	run(2, sample)
	run(10, sample)
	run(10, sample2)
}

func run(sz int, input string) {
	seen := map[pos]bool{}
	rope := make([]pos, sz)
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		dir, dist := parts[0], mustInt(parts[1])

		for dist > 0 {
			switch dir {
			case "L":
				rope[0] = rope[0].move(-1, 0)
			case "R":
				rope[0] = rope[0].move(1, 0)
			case "U":
				rope[0] = rope[0].move(0, -1)
			case "D":
				rope[0] = rope[0].move(0, 1)
			default:
				panic(dir)
			}
			for i := 1; i < sz; i++ {
				rope[i] = rope[i].moveTowards(rope[i-1])
			}
			seen[rope[sz-1]] = true
			dist--
		}
	}
	fmt.Println(len(seen))
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

func abs(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}
