package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

const sample = `
initial state: #..#.#..##......###...###

...## => #
..#.. => #
.#... => #
.#.#. => #
.#.## => #
.##.. => #
.#### => #
#.#.# => #
#.### => #
##.#. => #
##.## => #
###.. => #
###.# => #
####. => #
`

var (
	firstRe = regexp.MustCompile(`^initial state: ([.#]+)$`)
	otherRe = regexp.MustCompile(`^([.#]{5}) => ([.#])$`)
)

func main() {
	p := puz{
		states: map[int]bool{},
		min:    math.MaxInt,
		max:    math.MinInt,
	}

	first := true
	lines := strings.Split(sample, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if first {
			if !firstRe.MatchString(line) {
				panic(line)
			}
			m := firstRe.FindStringSubmatch(line)
			for i, c := range m[1] {
				switch c {
				case '#':
					p.states[i] = true
					p.min = min(p.min, i)
					p.max = max(p.max, i)
				case '.':
				default:
					panic(c)
				}
			}
			first = false
		} else {
			if !otherRe.MatchString(line) {
				panic(line)
			}
			m := otherRe.FindStringSubmatch(line)
			key := 0
			for _, c := range m[1] {
				key <<= 1
				switch c {
				case '#':
					key++
				case '.':
				default:
					panic(c)
				}
			}
			p.rules[key] = m[2][0] == '#'
		}
	}

	// part 1
	fmt.Println(string(p.Render()))
	cur := &p
	for i := 0; i < 20; i++ {
		cur = cur.Next()
		fmt.Println(string(cur.Render()))
	}
	fmt.Println(cur.Sum())

	// part 2
	const initTurns = 1000
	const finalTurns = 50000000000

	cur = &p
	last := cur.Sum()
	delta := 0
	for i := 0; i < initTurns; i++ {
		cur = cur.Next()
		next := cur.Sum()
		if i == initTurns/2 {
			delta = next - last
		} else if i > initTurns/2 {
			if next-last != delta {
				panic("inconsistent")
			}
		}
		last = next
	}
	project := delta*(finalTurns-initTurns) + last
	fmt.Printf("if sum=%d after %d turns, then sum=%d after %d turns\n", last, initTurns, project, finalTurns)
}

type puz struct {
	rules    [32]bool
	states   map[int]bool
	min, max int
}

func (p *puz) Next() *puz {
	newP := puz{
		rules:  p.rules,
		states: map[int]bool{},
		min:    math.MaxInt,
		max:    math.MinInt,
	}
	for i := p.min - 2; i <= p.max+2; i++ {
		key := p.Key(i)
		if newP.rules[key] {
			newP.states[i] = true
			newP.min = min(newP.min, i)
			newP.max = max(newP.max, i)
		}
	}

	return &newP
}

func (p *puz) Render() []byte {
	var line []byte
	for i := p.min; i <= p.max; i++ {
		if p.states[i] {
			line = append(line, '#')
		} else {
			line = append(line, '.')
		}
	}
	return line
}

func (p *puz) Size() int {
	return p.max - p.min
}

func (p *puz) Key(idx int) int {
	key := 0
	for i := idx - 2; i <= idx+2; i++ {
		key <<= 1
		if p.states[i] {
			key++
		}
	}
	return key
}

func (p *puz) Sum() int {
	sum := 0
	for i := p.min; i <= p.max; i++ {
		if p.states[i] {
			sum += i
		}
	}
	return sum
}
