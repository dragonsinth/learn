package main

import (
	"fmt"
	"strings"
)

var sample = `
5483143223
2745854711
5264556173
6141336146
6357385478
4167524645
2176841721
6882881134
4846848554
5283751526
`

var input = sample

type puzzle struct {
	state   [100]byte
	flashed [100]bool
}

func (p *puzzle) Step() int {
	p.flashed = [100]bool{}

	// Increment everything once
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			p.Inc(x, y)
		}
	}

	// Iteratively compute flashes.
	lastFlashes := 0
	for {
		for y := 0; y < 10; y++ {
			for x := 0; x < 10; x++ {
				if p.state[y*10+x] == 10 {
					p.state[y*10+x] = 0
					p.flashed[y*10+x] = true

					// maybe flash neighbors
					for j := y - 1; j <= y+1; j++ {
						for i := x - 1; i <= x+1; i++ {
							p.Spill(i, j)
						}
					}
				}
			}
		}

		nextFlashes := 0
		for _, flashed := range p.flashed {
			if flashed {
				nextFlashes++
			}
		}
		if nextFlashes == lastFlashes {
			return nextFlashes
		}
		lastFlashes = nextFlashes
	}
}

func (p *puzzle) Inc(x, y int) {
	if x < 0 || x >= 10 {
		return
	}
	if y < 0 || y >= 10 {
		return
	}
	if p.state[y*10+x] < 10 {
		p.state[y*10+x]++
	}
}

func (p *puzzle) Spill(x, y int) {
	if x < 0 || x >= 10 {
		return
	}
	if y < 0 || y >= 10 {
		return
	}
	if p.flashed[y*10+x] {
		return
	}
	if p.state[y*10+x] < 10 {
		p.state[y*10+x]++
	}
}

func (p *puzzle) Print() {
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			fmt.Printf("%c", p.state[y*10+x]+'0')
		}
		fmt.Println()
	}
}

func main() {
	var p puzzle
	y := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		buf := []byte(line)
		if len(buf) != 10 {
			panic(buf)
		}

		for x, c := range buf {
			p.state[y*10+x] = c - '0'
		}
		y++
	}
	if y != 10 {
		panic(y)
	}

	p.Print()

	sum := 0
	for step := 1; ; step++ {
		fmt.Println("Step:", step)
		score := p.Step()
		p.Print()
		fmt.Println()

		if step <= 100 {
			sum += score
		} else if score == 100 {
			break
		}
	}
	fmt.Println(sum)
}
