package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var sample = `
6,10
0,14
9,10
0,3
10,4
4,11
6,0
6,12
4,1
0,13
10,12
3,4
3,0
8,4
1,10
2,14
8,10
9,0

fold along y=7
fold along x=5
`

var input = sample

const debug = true

type puzzle struct {
	on     [][]bool
	height int
	width  int
}

func (p *puzzle) Print(debug bool) {
	fmt.Println()
	fmt.Println(p.width, p.height)
	sum := 0
	for _, row := range p.on {
		for _, v := range row {
			if v {
				sum++
			}
			if debug {
				if v {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
		}
		if debug {
			fmt.Println()
		}
	}
	fmt.Println(sum, "dots")
}

func (p *puzzle) FoldHorizontal(foldY int) {
	// Fold rows together.
	yEnd := 2*foldY + 1
	for y := 0; y < foldY; y++ {
		yComp := yEnd - y - 1
		if yComp >= p.height {
			continue
		}
		row1 := p.on[y]
		row2 := p.on[yComp]
		for x := range row1 {
			row1[x] = row1[x] || row2[x]
		}
	}

	// Trim.
	p.height = foldY
	p.on = p.on[:foldY]
}

func (p *puzzle) FoldVertically(x int) {
	expect := (p.width - 1) / 2
	if x != expect {
		panic(fmt.Sprint("invalid fold", x, expect))
	}

	// Fold each row onto itself, then trim.
	for y, row := range p.on {
		for x := 0; x < expect; x++ {
			row[x] = row[x] || row[p.width-x-1]
			p.on[y] = row[:expect]
		}
	}

	p.width = expect
}

type fold struct {
	x, y int
}

var foldRegex = regexp.MustCompile(`^fold along (x|y)=(\d+)$`)

func main() {
	var puz puzzle
	var folds []fold

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "fold along ") {
			submatches := foldRegex.FindStringSubmatch(line)
			switch submatches[1] {
			case "x":
				folds = append(folds, fold{
					x: mustInt(submatches[2]),
				})
			case "y":
				folds = append(folds, fold{
					y: mustInt(submatches[2]),
				})
			default:
				panic(submatches[1])
			}
			continue
		}

		parts := strings.Split(line, ",")
		x, y := mustInt(parts[0]), mustInt(parts[1])
		if x >= puz.width {
			puz.width = x + 1
		}
		if y >= puz.height {
			puz.height = y + 1
		}
	}

	fmt.Println(folds)
	for y := 0; y < puz.height; y++ {
		puz.on = append(puz.on, make([]bool, puz.width))
	}

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "fold along ") {
			continue
		}

		parts := strings.Split(line, ",")
		x, y := mustInt(parts[0]), mustInt(parts[1])
		puz.on[y][x] = true
	}

	puz.Print(debug)
	for _, f := range folds {
		switch {
		case f.x > 0:
			puz.FoldVertically(f.x)
		case f.y > 0:
			puz.FoldHorizontal(f.y)
		default:
			panic(f)
		}
		puz.Print(debug)
	}
	puz.Print(true)
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
