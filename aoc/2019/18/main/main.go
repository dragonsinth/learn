package main

import (
	"container/heap"
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"os"
	"sort"
	"strings"
)

var samples = []string{
	`
#########
#b.A.@.a#
#########`,
	`
########################
#f.D.E.e.C.b.A.@.a.B.c.#
######################.#
#d.....................#
########################`,
	`
########################
#...............b.C.D.f#
#.######################
#.....@.a.B.c.d.A.e.F.g#
########################`,
	`
#################
#i.G..c...e..H.p#
########.########
#j.A..b...f..D.o#
########@########
#k.E..a...g..B.n#
########.########
#l.F..d...h..C.m#
#################`,
	`
########################
#@..............ac.GI.b#
###d#e#f################
###A#B#C################
###g#h#i################
########################`,
}

var samples2 = []string{
	`
#######
#a.#Cd#
##@#@##
#######
##@#@##
#cB#Ab#
#######`,
	`
###############
#d.ABC.#.....a#
######@#@######
###############
######@#@######
#b.....#.....c#
###############`,
	`
#############
#DcBa.#.GhKl#
#.###@#@#I###
#e#d#####j#k#
###C#@#@###J#
#fEbA.#.FgHi#
#############`,
	`
#############
#g#f.D#..h#l#
#F###e#E###.#
#dCba@#@BcIJ#
#############
#nK.L@#@G...#
#M###N#H###.#
#o#m..#i#jk.#
#############`,
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

type dir int

const (
	N = dir(iota)
	E
	S
	W
	MAX_DIR
)

type key byte // a-z

func (k key) door() door {
	return door(k + 'A' - 'a')
}

type door byte // A-Z

func (d door) key() key {
	return key(d + 'a' - 'A')
}

func main() {
	for i, input := range append(samples, samples2...) {
		fmt.Printf("%d\n", i)
		run(input)
		fmt.Println()
	}
}

func run(input string) {
	p := &puzzle{
		megaMap: map[pt][]edge{},
	}

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		p.maze = append(p.maze, []byte(line))
		if len(p.maze[0]) != len(line) {
			panic(line)
		}
	}

	// Find players and keys.
	var players []pt
	keys := map[key]pt{}
	for y, line := range p.maze {
		for x, v := range line {
			pos := pt{x, y}
			if v == '@' {
				players = append(players, pos)
			}
			if v >= 'a' && v <= 'z' {
				keys[key(v)] = pos
			}
		}
	}

	termbox.RenderPlain(p.maze, os.Stdout)

	// Create a giant graph of nodes and edges.
	for _, v := range keys {
		edges := p.findAllKeys(v)
		p.megaMap[v] = edges
	}
	for _, v := range players {
		edges := p.findAllKeys(v)
		p.megaMap[v] = edges
	}

	fmt.Println(p.Solve(players, len(keys)))
}

func (p *puzzle) findAllKeys(pos pt) []edge {
	type loc struct {
		pos   pt
		dist  int
		doors string
	}

	start := loc{
		pos:   pos,
		dist:  0,
		doors: "",
	}
	locMap := map[pt]loc{pos: start}
	work := []loc{start}

	var ret []edge
	for i := 0; i < len(work); i++ {
		w := work[i]
		_, ok := locMap[w.pos]
		if !ok {
			panic("wat")
		}
		if k, ok := p.isKey(w.pos); ok {
			ret = append(ret, edge{
				a:     pos,
				b:     w.pos,
				k:     k,
				dist:  w.dist,
				doors: w.doors,
			})
		}

		for dir := N; dir < MAX_DIR; dir++ {
			dst := w.pos.move(dir)
			if _, ok := locMap[dst]; ok {
				continue // already computed
			}
			if p.isWall(dst) {
				continue
			}

			newLoc := loc{
				pos:   dst,
				dist:  w.dist + 1,
				doors: w.doors,
			}

			// If there's a door here
			if d, ok := p.isDoor(dst); ok {
				newLoc.doors = w.doors + string(d)
			}

			locMap[dst] = newLoc
			work = append(work, newLoc)
		}
	}
	return ret
}

type edge struct {
	a, b  pt
	k     key
	dist  int
	doors string
}

type puzzle struct {
	maze    [][]byte
	megaMap map[pt][]edge
}

type partialPath struct {
	dist int
	pos  []pt
	path string
}

func (p *puzzle) Solve(players []pt, nKeys int) (int, string) {
	h := workHeap{{pos: players}}
	heap.Init(&h)

	type scoredPosition struct {
		sortedKeys string
		pos        [4]pt
	}
	scoredKeySets := map[scoredPosition]int{}

	for {
		pp := heap.Pop(&h).(partialPath)
		if len(pp.path) == nKeys {
			return pp.dist, pp.path
		}

		sp := scoredPosition{
			sortedKeys: sortedString(pp.path),
		}
		copy(sp.pos[:], pp.pos)

		if v, ok := scoredKeySets[sp]; ok && pp.dist >= v {
			continue
		}
		scoredKeySets[sp] = pp.dist

		for i, pos := range pp.pos {
			found := p.findAccessibleKeys(pos, pp.path)
			for _, e := range found {
				newPos := append([]pt{}, pp.pos...)
				newPos[i] = e.pos
				newPP := partialPath{
					dist: pp.dist + e.dist,
					pos:  newPos,
					path: pp.path + string(e.k),
				}
				heap.Push(&h, newPP)
			}
		}
	}
}

func sortedString(in string) string {
	ret := []byte(in)
	sort.Slice(ret, func(i, j int) bool {
		return ret[i] < ret[j]
	})
	return string(ret)
}

type keyLoc struct {
	k    key
	pos  pt
	dist int
}

func (p *puzzle) findAccessibleKeys(from pt, path string) []keyLoc {
	edges, ok := p.megaMap[from]
	if !ok {
		panic(from)
	}

	var ret []keyLoc
	for _, e := range edges {
		if haveKey(e.k, path) {
			// already have this one
			continue
		}

		// do we have all the keys?
		haveKeys := 0
		for _, d := range e.doors {
			if haveKey(door(d).key(), path) {
				haveKeys++
			}
		}

		if haveKeys == len(e.doors) {
			ret = append(ret, keyLoc{
				k:    e.k,
				pos:  e.b,
				dist: e.dist,
			})
		}
	}
	return ret
}

func haveKey(k key, path string) bool {
	for _, e := range path {
		if k == key(e) {
			return true
		}
	}
	return false
}

func (p *puzzle) isKey(w pt) (key, bool) {
	v := p.maze[w.Y][w.X]
	if v >= 'a' && v <= 'z' {
		return key(v), true
	}
	return 0, false
}

func (p *puzzle) isDoor(w pt) (door, bool) {
	v := p.maze[w.Y][w.X]
	if v >= 'A' && v <= 'Z' {
		return door(v), true
	}
	return 0, false
}

func (p *puzzle) isWall(w pt) bool {
	if p.maze[w.Y][w.X] == '#' {
		return true
	}
	return false
}

type workHeap []partialPath

func (w *workHeap) Len() int {
	return len(*w)
}

func (w *workHeap) Less(i, j int) bool {
	return (*w)[i].dist < (*w)[j].dist
}

func (w *workHeap) Swap(i, j int) {
	(*w)[i], (*w)[j] = (*w)[j], (*w)[i]
}

func (w *workHeap) Push(x any) {
	*w = append(*w, x.(partialPath))
}

func (w *workHeap) Pop() any {
	x := (*w)[len(*w)-1]
	*w = (*w)[0 : len(*w)-1]
	return x
}
