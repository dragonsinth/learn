package main

import (
	"fmt"
	"strings"
)

var sample = `
1 + (2 * 3) + (4 * (5 + 6))
2 * 3 + (4 * 5)
5 + (8 * 3 + 9 + 3 * 4 * 3)
5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))
((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2
`

var input = sample

func main() {
	sum := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		val := parse(line)
		fmt.Println(val)
		sum += val
	}

	fmt.Println("----")
	fmt.Println(sum)
}

type stream struct {
	line string
	pos  int
}

func (s *stream) Peek() byte {
	return s.line[s.pos]
}

func (s *stream) Pop() byte {
	c := s.line[s.pos]
	s.pos++
	return c
}

func (s *stream) MustPop(expect byte) {
	c := s.Pop()
	if c != expect {
		panic(fmt.Sprintf("expected %s, got %s", string(expect), string(c)))
	}
}

func (s *stream) Eof() bool {
	return s.pos >= len(s.line)
}
