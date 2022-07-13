package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"strings"
)

var samples = []string{
	`
         A           
         A           
  #######.#########  
  #######.........#  
  #######.#######.#  
  #######.#######.#  
  #######.#######.#  
  #####  B    ###.#  
BC...##  C    ###.#  
  ##.##       ###.#  
  ##...DE  F  ###.#  
  #####    G  ###.#  
  #########.#####.#  
DE..#######...###.#  
  #.#########.###.#  
FG..#########.....#  
  ###########.#####  
             Z       
             Z       `,
	`
                   A               
                   A               
  #################.#############  
  #.#...#...................#.#.#  
  #.#.#.###.###.###.#########.#.#  
  #.#.#.......#...#.....#.#.#...#  
  #.#########.###.#####.#.#.###.#  
  #.............#.#.....#.......#  
  ###.###########.###.#####.#.#.#  
  #.....#        A   C    #.#.#.#  
  #######        S   P    #####.#  
  #.#...#                 #......VT
  #.#.#.#                 #.#####  
  #...#.#               YN....#.#  
  #.###.#                 #####.#  
DI....#.#                 #.....#  
  #####.#                 #.###.#  
ZZ......#               QG....#..AS
  ###.###                 #######  
JO..#.#.#                 #.....#  
  #.#.#.#                 ###.#.#  
  #...#..DI             BU....#..LF
  #####.#                 #.#####  
YN......#               VT..#....QG
  #.###.#                 #.###.#  
  #.#...#                 #.....#  
  ###.###    J L     J    #.#.###  
  #.....#    O F     P    #.#...#  
  #.###.#####.#.#####.#####.###.#  
  #...#.#.#...#.....#.....#.#...#  
  #.#####.###.###.#.#.#########.#  
  #...#.#.....#...#.#.#.#.....#.#  
  #.###.#####.###.###.#.#.#######  
  #.#.........#...#.............#  
  #########.###.###.#############  
           B   J   C               
           U   P   P               `,
	`
             Z L X W       C                 
             Z P Q B       K                 
  ###########.#.#.#.#######.###############  
  #...#.......#.#.......#.#.......#.#.#...#  
  ###.#.#.#.#.#.#.#.###.#.#.#######.#.#.###  
  #.#...#.#.#...#.#.#...#...#...#.#.......#  
  #.###.#######.###.###.#.###.###.#.#######  
  #...#.......#.#...#...#.............#...#  
  #.#########.#######.#.#######.#######.###  
  #...#.#    F       R I       Z    #.#.#.#  
  #.###.#    D       E C       H    #.#.#.#  
  #.#...#                           #...#.#  
  #.###.#                           #.###.#  
  #.#....OA                       WB..#.#..ZH
  #.###.#                           #.#.#.#  
CJ......#                           #.....#  
  #######                           #######  
  #.#....CK                         #......IC
  #.###.#                           #.###.#  
  #.....#                           #...#.#  
  ###.###                           #.#.#.#  
XF....#.#                         RF..#.#.#  
  #####.#                           #######  
  #......CJ                       NM..#...#  
  ###.#.#                           #.###.#  
RE....#.#                           #......RF
  ###.###        X   X       L      #.#.#.#  
  #.....#        F   Q       P      #.#.#.#  
  ###.###########.###.#######.#########.###  
  #.....#...#.....#.......#...#.....#.#...#  
  #####.#.###.#######.#######.###.###.#.#.#  
  #.......#.......#.#.#.#.#...#...#...#.#.#  
  #####.###.#####.#.#.#.#.###.###.#.###.###  
  #.......#.....#.#...#...............#...#  
  #############.#.#.###.###################  
               A O F   N                     
               A A D   M                     `,
}

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

func (p pt) z(z int) pt3 {
	return pt3{p.X, p.Y, z}
}

type dir int

func (d dir) opposite() dir {
	return (d + 2) % 4
}

type edge struct {
	inner pt
	outer pt
}

const (
	N = dir(iota)
	E
	S
	W
	MAX_DIR
)

func main() {
	for _, sample := range samples {
		run(sample)
	}
}

func run(input string) {
	var data [][]byte
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		data = append(data, []byte(line))
	}

	portals := map[pt]string{}
	edges := map[string]edge{}
	maze := map[pt]bool{}

	// Map all the portals
	for y, line := range data {
		for x, v := range line {
			if x == 0 || y == 0 || x == len(line)-1 || y == len(data)-1 {
				continue // edge of map
			}
			pos := pt{x, y}
			if v == '.' {
				maze[pos] = true
			}
			if v >= 'A' && v <= 'Z' {
				// Find the adjacent dot.
				for d := N; d < MAX_DIR; d++ {
					dst := pos.move(d)
					if data[dst.Y][dst.X] == '.' {
						// found it
						opp := pos.move(d.opposite())
						wat := data[opp.Y][opp.X]
						var name string
						switch d {
						case N, W:
							name = string([]byte{v, wat})
						case S, E:
							name = string([]byte{wat, v})
						}
						portals[dst] = name
						e := edges[name]
						if x == 1 || y == 1 || x == len(line)-2 || y == len(data)-2 {
							e.outer = dst
						} else {
							e.inner = dst
						}
						edges[name] = e
					}
				}
			}
		}
	}

	termbox.RenderPlain(data, os.Stdout)
	distMap := computeDistMap(maze, edges["AA"].outer, portals, edges)
	fmt.Println("2d:", distMap[edges["ZZ"].outer])
	distMap3 := computeDistMap3(maze, edges["AA"].outer.z(0), portals, edges)
	fmt.Println("3d:", distMap3[edges["ZZ"].outer.z(0)])
}

func computeDistMap(maze map[pt]bool, start pt, portals map[pt]string, edges map[string]edge) map[pt]int {
	distMap := map[pt]int{
		start: 0,
	}
	work := []pt{start}
	for i := 0; i < len(work); i++ {
		w := work[i]
		dist, ok := distMap[w]
		if !ok {
			panic("wat")
		}
		var checkPoints []pt
		for dir := N; dir < MAX_DIR; dir++ {
			checkPoints = append(checkPoints, w.move(dir))
		}

		// See if we're a portal connected to the other side!
		if name, ok := portals[w]; ok {
			if name == "ZZ" {
				return distMap // early exit
			}
			checkPoints = append(checkPoints, edges[name].inner)
			checkPoints = append(checkPoints, edges[name].outer)
		}

		for _, dst := range checkPoints {
			if _, ok := distMap[dst]; ok {
				continue // already computed
			}
			if maze[dst] {
				distMap[dst] = dist + 1
				work = append(work, dst)
			}
		}

	}
	return distMap
}

type pt3 struct {
	X, Y, Z int
}

func (p pt3) move(d dir) pt3 {
	switch d {
	case N:
		return pt3{p.X, p.Y - 1, p.Z}
	case E:
		return pt3{p.X + 1, p.Y, p.Z}
	case S:
		return pt3{p.X, p.Y + 1, p.Z}
	case W:
		return pt3{p.X - 1, p.Y, p.Z}
	default:
		panic(d)
	}
}

func (p pt3) flat() pt {
	return pt{p.X, p.Y}
}

func computeDistMap3(maze map[pt]bool, start pt3, portals map[pt]string, edges map[string]edge) map[pt3]int {
	distMap := map[pt3]int{
		start: 0,
	}
	work := []pt3{start}
	for i := 0; i < len(work); i++ {
		w := work[i]
		dist, ok := distMap[w]
		if !ok {
			panic("wat")
		}
		if dist > 9999 {
			continue // probably no solution?
		}

		var checkPoints []pt3
		for dir := N; dir < MAX_DIR; dir++ {
			checkPoints = append(checkPoints, w.move(dir))
		}

		// See if we're a portal connected to the other side!
		if name, ok := portals[w.flat()]; ok {
			if w.Z == 0 {
				if name == "ZZ" {
					return distMap // early exit
				}
			} else if name == "AA" || name == "ZZ" {
				continue // this is actually a wall
			}

			e := edges[name]
			if e.inner == w.flat() {
				// traverse inner -> outer goes down a level
				checkPoints = append(checkPoints, e.outer.z(w.Z+1))
			} else if e.outer == w.flat() {
				// traverse outer -> inner goes up a level unless at 0
				if w.Z > 0 {
					checkPoints = append(checkPoints, e.inner.z(w.Z-1))
				}
			} else {
				panic(e)
			}
		}

		for _, dst := range checkPoints {
			if _, ok := distMap[dst]; ok {
				continue // already computed
			}
			if maze[dst.flat()] {
				distMap[dst] = dist + 1
				work = append(work, dst)
			}
		}

	}
	return distMap
}
