package main

import (
	"fmt"
	"regexp"
	"strings"
)

var sample = `
.#.
..#
###
`

var input = sample

var parseLine = regexp.MustCompile(`^[.#]+$`)

func main() {
	width := 0
	var zeroPlane [][]byte
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !parseLine.MatchString(line) {
			panic(line)
		}

		buf := []byte(line)
		if width == 0 {
			width = len(buf)
		} else if width != len(buf) {
			panic(line)
		}

		row := make([]byte, width)
		for i, c := range buf {
			if c == '#' {
				row[i] = 1
			}
		}
		zeroPlane = append(zeroPlane, row)
	}

	puz := &puzzle2{
		v:      [][][][]byte{{zeroPlane}},
		width:  width,
		height: len(zeroPlane),
		depth:  1,
		zonk:   1,
	}
	puz.Print()

	cur := puz
	for i := 0; i < 6; i++ {
		cur = cur.Next()
		if i < 2 {
			cur.Print()
		} else {
			fmt.Println(cur.Active())
		}
	}
}
