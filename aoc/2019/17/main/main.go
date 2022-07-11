package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/2019/intcode"
	"github.com/dragonsinth/learn/aoc/termbox"
	"strconv"
	"strings"
	"time"
)

const data = `1,330,331,332,109,3028,1102,1182,1,16,1101,0,1437,24,101,0,0,570,1006,570,36,102,1,571,0,1001,570,-1,570,1001,24,1,24,1105,1,18,1008,571,0,571,1001,16,1,16,1008,16,1437,570,1006,570,14,21101,58,0,0,1106,0,786,1006,332,62,99,21101,0,333,1,21102,1,73,0,1105,1,579,1101,0,0,572,1101,0,0,573,3,574,101,1,573,573,1007,574,65,570,1005,570,151,107,67,574,570,1005,570,151,1001,574,-64,574,1002,574,-1,574,1001,572,1,572,1007,572,11,570,1006,570,165,101,1182,572,127,101,0,574,0,3,574,101,1,573,573,1008,574,10,570,1005,570,189,1008,574,44,570,1006,570,158,1106,0,81,21102,340,1,1,1105,1,177,21101,477,0,1,1105,1,177,21101,0,514,1,21102,176,1,0,1105,1,579,99,21102,184,1,0,1105,1,579,4,574,104,10,99,1007,573,22,570,1006,570,165,101,0,572,1182,21102,375,1,1,21101,211,0,0,1106,0,579,21101,1182,11,1,21102,1,222,0,1105,1,979,21101,388,0,1,21102,233,1,0,1106,0,579,21101,1182,22,1,21102,244,1,0,1106,0,979,21101,0,401,1,21102,1,255,0,1106,0,579,21101,1182,33,1,21101,266,0,0,1106,0,979,21101,414,0,1,21102,277,1,0,1105,1,579,3,575,1008,575,89,570,1008,575,121,575,1,575,570,575,3,574,1008,574,10,570,1006,570,291,104,10,21102,1182,1,1,21102,313,1,0,1105,1,622,1005,575,327,1101,1,0,575,21102,327,1,0,1105,1,786,4,438,99,0,1,1,6,77,97,105,110,58,10,33,10,69,120,112,101,99,116,101,100,32,102,117,110,99,116,105,111,110,32,110,97,109,101,32,98,117,116,32,103,111,116,58,32,0,12,70,117,110,99,116,105,111,110,32,65,58,10,12,70,117,110,99,116,105,111,110,32,66,58,10,12,70,117,110,99,116,105,111,110,32,67,58,10,23,67,111,110,116,105,110,117,111,117,115,32,118,105,100,101,111,32,102,101,101,100,63,10,0,37,10,69,120,112,101,99,116,101,100,32,82,44,32,76,44,32,111,114,32,100,105,115,116,97,110,99,101,32,98,117,116,32,103,111,116,58,32,36,10,69,120,112,101,99,116,101,100,32,99,111,109,109,97,32,111,114,32,110,101,119,108,105,110,101,32,98,117,116,32,103,111,116,58,32,43,10,68,101,102,105,110,105,116,105,111,110,115,32,109,97,121,32,98,101,32,97,116,32,109,111,115,116,32,50,48,32,99,104,97,114,97,99,116,101,114,115,33,10,94,62,118,60,0,1,0,-1,-1,0,1,0,0,0,0,0,0,1,8,18,0,109,4,2101,0,-3,586,21002,0,1,-1,22101,1,-3,-3,21101,0,0,-2,2208,-2,-1,570,1005,570,617,2201,-3,-2,609,4,0,21201,-2,1,-2,1105,1,597,109,-4,2106,0,0,109,5,2101,0,-4,629,21002,0,1,-2,22101,1,-4,-4,21102,1,0,-3,2208,-3,-2,570,1005,570,781,2201,-4,-3,652,21002,0,1,-1,1208,-1,-4,570,1005,570,709,1208,-1,-5,570,1005,570,734,1207,-1,0,570,1005,570,759,1206,-1,774,1001,578,562,684,1,0,576,576,1001,578,566,692,1,0,577,577,21101,0,702,0,1105,1,786,21201,-1,-1,-1,1106,0,676,1001,578,1,578,1008,578,4,570,1006,570,724,1001,578,-4,578,21102,1,731,0,1106,0,786,1105,1,774,1001,578,-1,578,1008,578,-1,570,1006,570,749,1001,578,4,578,21101,0,756,0,1105,1,786,1105,1,774,21202,-1,-11,1,22101,1182,1,1,21102,1,774,0,1105,1,622,21201,-3,1,-3,1105,1,640,109,-5,2105,1,0,109,7,1005,575,802,21002,576,1,-6,20102,1,577,-5,1106,0,814,21102,0,1,-1,21101,0,0,-5,21101,0,0,-6,20208,-6,576,-2,208,-5,577,570,22002,570,-2,-2,21202,-5,37,-3,22201,-6,-3,-3,22101,1437,-3,-3,2101,0,-3,843,1005,0,863,21202,-2,42,-4,22101,46,-4,-4,1206,-2,924,21101,0,1,-1,1105,1,924,1205,-2,873,21102,1,35,-4,1105,1,924,1201,-3,0,878,1008,0,1,570,1006,570,916,1001,374,1,374,1201,-3,0,895,1102,1,2,0,1201,-3,0,902,1001,438,0,438,2202,-6,-5,570,1,570,374,570,1,570,438,438,1001,578,558,921,21002,0,1,-4,1006,575,959,204,-4,22101,1,-6,-6,1208,-6,37,570,1006,570,814,104,10,22101,1,-5,-5,1208,-5,43,570,1006,570,810,104,10,1206,-1,974,99,1206,-1,974,1102,1,1,575,21101,973,0,0,1106,0,786,99,109,-7,2106,0,0,109,6,21101,0,0,-4,21101,0,0,-3,203,-2,22101,1,-3,-3,21208,-2,82,-1,1205,-1,1030,21208,-2,76,-1,1205,-1,1037,21207,-2,48,-1,1205,-1,1124,22107,57,-2,-1,1205,-1,1124,21201,-2,-48,-2,1106,0,1041,21101,-4,0,-2,1106,0,1041,21101,0,-5,-2,21201,-4,1,-4,21207,-4,11,-1,1206,-1,1138,2201,-5,-4,1059,1202,-2,1,0,203,-2,22101,1,-3,-3,21207,-2,48,-1,1205,-1,1107,22107,57,-2,-1,1205,-1,1107,21201,-2,-48,-2,2201,-5,-4,1090,20102,10,0,-1,22201,-2,-1,-2,2201,-5,-4,1103,2102,1,-2,0,1106,0,1060,21208,-2,10,-1,1205,-1,1162,21208,-2,44,-1,1206,-1,1131,1106,0,989,21101,0,439,1,1106,0,1150,21101,477,0,1,1105,1,1150,21102,514,1,1,21102,1149,1,0,1106,0,579,99,21102,1157,1,0,1105,1,579,204,-2,104,10,99,21207,-3,22,-1,1206,-1,1138,1201,-5,0,1176,1201,-4,0,0,109,-6,2106,0,0,26,7,30,1,5,1,30,1,1,9,26,1,1,1,3,1,3,1,26,1,1,1,3,1,3,1,26,1,1,1,3,1,3,1,26,5,1,1,3,1,28,1,1,1,1,1,3,1,28,1,1,1,1,1,3,1,28,1,1,1,1,1,3,1,28,5,3,1,30,1,5,1,28,9,28,1,1,1,28,9,34,1,36,1,36,1,12,5,11,9,12,1,15,1,20,1,15,1,20,1,15,1,20,1,15,1,20,1,15,1,16,11,5,11,10,1,3,1,5,1,5,1,3,1,5,1,10,1,3,1,1,9,1,1,3,5,1,1,10,1,3,1,1,1,3,1,3,1,1,1,7,1,1,1,10,1,3,7,3,1,1,1,7,1,1,1,10,1,5,1,7,1,1,1,7,1,1,1,10,7,7,1,1,1,7,1,1,1,24,1,1,1,7,1,1,1,24,1,1,9,1,7,18,1,17,1,18,7,1,9,1,1,24,1,1,1,7,1,1,1,24,1,1,1,7,1,1,1,24,1,1,1,7,1,1,1,24,1,1,11,24,1,9,1,26,1,9,1,26,1,9,1,26,11,6`

type pt struct {
	X, Y int
}

func (p pt) move(d dir) pt {
	switch d {
	case N:
		return pt{p.X, p.Y - 1}
	case E:
		return pt{p.X + 1, p.Y}
	case S:
		return pt{p.X, p.Y + 1}
	case W:
		return pt{p.X - 1, p.Y}
	default:
		panic(d)
	}
}

type dir int

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

const (
	w, h = 37, 45

	prog = `A,A,B,C,A,C,B,C,A,B
L,4,L,10,L,6
L,6,L,4,R,8,R,8
L,6,R,8,L,10,L,8,L,8
y
`
)

func main() {
	main2()
}

func main1() {
	codes := intcode.Parse(data)

	var p puzzle
	for y := 0; y < h; y++ {
		p.data = append(p.data, make([]byte, w))
	}

	var cursor pt
	m := intcode.NewIntMachine(codes, nil, func(v int) {
		if v == '\n' {
			cursor.Y++
			cursor.X = 0
		} else {
			p.data[cursor.Y][cursor.X] = byte(v)
			cursor.X++
		}
	})
	m.Run()

	var bot, tgt pt
	var direction dir

	tot := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			loc := pt{x, y}
			switch p.data[y][x] {
			case '^':
				bot = loc
				direction = N
			case '>':
				bot = loc
				direction = E
			case 'v':
				bot = loc
				direction = S
			case '<':
				bot = loc
				direction = W
			case '#':
				sum := p.IsPath(loc.move(N)) +
					p.IsPath(loc.move(S)) +
					p.IsPath(loc.move(E)) +
					p.IsPath(loc.move(W))
				switch sum {
				case 1:
					tgt = loc
				case 4:
					tot += loc.X * loc.Y
				}

			}
		}
	}
	for _, line := range p.data {
		fmt.Println(string(line))
	}
	fmt.Println(tot)

	fmt.Println("bot: ", bot)
	fmt.Println("tgt: ", tgt)
	// Compute a path from bot to tgt.
	var steps []string
	for bot != tgt {
		// Find direction
		if p.IsPath(bot.move(direction.left())) > 0 {
			direction = direction.left()
			steps = append(steps, "L")
		} else if p.IsPath(bot.move(direction.right())) > 0 {
			direction = direction.right()
			steps = append(steps, "R")
		} else {
			fmt.Println(steps)
			panic(bot)
		}

		// Walk forward
		dist := 0
		for p.IsPath(bot.move(direction)) > 0 {
			bot = bot.move(direction)
			dist++
		}
		steps = append(steps, strconv.Itoa(dist))
	}
	fmt.Println(strings.Join(steps, ","))
}

func main2() {
	codes := intcode.Parse(data)
	term := termbox.New(false)
	defer term.Stop()

	var p puzzle
	for y := 0; y < h; y++ {
		p.data = append(p.data, make([]byte, w))
	}

	codes[0] = 2
	inputIdx := 0
	var line []byte
	var cursor pt
	rendering := false
	m := intcode.NewIntMachine(codes, func() int {
		r := int(prog[inputIdx])
		inputIdx++
		if inputIdx == len(prog) {
			term.Start()
			rendering = true
		}
		return r
	}, func(v int) {
		if v > 255 {
			term.Stop()
			fmt.Println(v)
			return
		}

		if rendering {
			// map mode
			if v == '\n' {
				if cursor.X == w {
					cursor.Y++
					cursor.X = 0
				} else {
					if term.Enabled() {
						time.Sleep(10 * time.Millisecond)
						term.Render(p.data, nil)
					}
					cursor.Y = 0
				}
			} else {
				p.data[cursor.Y][cursor.X] = byte(v)
				cursor.X++
			}
		} else {
			// line by line mode
			if v == '\n' {
				fmt.Println(string(line))
				line = nil
			} else {
				line = append(line, byte(v))
			}
		}
	})
	m.Run()
}

type puzzle struct {
	data [][]byte
}

func (p puzzle) IsPath(loc pt) int {
	if loc.X < 0 || loc.Y < 0 || loc.X >= w || loc.Y >= h {
		return 0
	}
	switch p.data[loc.Y][loc.X] {
	case '#', '^', '>', 'v', '<':
		return 1
	default:
		return 0
	}
}
