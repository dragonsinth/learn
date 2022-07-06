package main

import (
	"fmt"
	"math"
)

var sample = `0222112222120000`

func main() {
	im := parse(sample, 2, 2)

	best := []int{math.MaxInt32}
	for _, l := range im.layers {
		counts := l.counts()
		if counts[0] < best[0] {
			best = counts
		}
	}
	fmt.Println(best, best[1]*best[2])
	fmt.Println(im.render())
}

type image struct {
	w, h   int
	layers []layer
}

func (i *image) render() string {
	var d [][]byte
	for y := 0; y < i.h; y++ {
		line := make([]byte, i.w)
		for x := range line {
			line[x] = 2
		}
		d = append(d, line)
	}
	for _, l := range i.layers {
		l.render(d)
	}
	var out []byte
	for _, line := range d {
		for _, v := range line {
			if v == 1 {
				out = append(out, '*')
			} else {
				out = append(out, ' ')
			}
		}
		out = append(out, '\n')
	}
	return string(out)
}

type layer struct {
	data [][]byte
}

func (l layer) counts() []int {
	ret := make([]int, 10)
	for _, line := range l.data {
		for _, v := range line {
			ret[v]++
		}
	}
	return ret
}

func (l layer) render(d [][]byte) {
	for y, line := range d {
		for x, v := range line {
			if v == 2 {
				d[y][x] = l.data[y][x]
			}
		}
	}
}

func parse(input string, w int, h int) *image {
	ret := &image{
		w:      w,
		h:      h,
		layers: nil,
	}
	i := 0
	for i < len(input) {
		var l layer
		for y := 0; y < h; y++ {
			line := make([]byte, w)
			for x := range line {
				line[x] = input[i] - '0'
				i++
			}
			l.data = append(l.data, line)
		}
		ret.layers = append(ret.layers, l)
	}

	return ret
}
