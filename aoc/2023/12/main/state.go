package main

import "fmt"

type state struct {
	first     byte
	rest      string
	acc       byte
	expect    [32]byte
	expectPos int
	expectLen int
}

func (s state) Count(seen map[state]int) int {
	if s.terminal() {
		return 1
	}
	if v, ok := seen[s]; ok {
		return v
	}

	sum := 0
	for _, sp := range s.nextValid() {
		sum += sp.Count(seen)
	}
	seen[s] = sum
	return sum
}

func (s state) String() string {
	return fmt.Sprintf("%c%s: %d %+v", s.first, s.rest, s.acc, s.expect[s.expectPos:s.expectLen])
}

func (s state) terminal() bool {
	if s.first == 0 {
		if s.rest != "" {
			panic(s.rest)
		}
		if s.expectPos != s.expectLen {
			panic(s.expectPos)
		}
		if s.acc != 0 {
			panic(s.acc)
		}
		return true
	}
	return false
}

func (s state) nextValid() []state {
	if s.terminal() {
		panic("here")
	}

	checkTerminal := func(s state) []state {
		s, ok := s.maybeAccumulate()
		if !ok {
			return nil
		}
		if s.expectPos != s.expectLen {
			return nil // did not consume everything
		}
		s.first = 0
		return []state{s} // ok
	}

	switch s.first {
	case '#':
		if s.expectPos >= s.expectLen {
			return nil // no more valid states
		}
		s.acc++
		if s.acc > s.expect[s.expectPos] {
			return nil // overflow
		}

		if s.rest == "" {
			return checkTerminal(s)
		} else {
			s.first = s.rest[0]
			s.rest = s.rest[1:]
		}

		return []state{s}
	case '.':
		s, ok := s.maybeAccumulate()
		if !ok {
			return nil
		}

		if s.rest == "" {
			return checkTerminal(s)
		} else {
			s.first = s.rest[0]
			s.rest = s.rest[1:]
		}

		return []state{s}
	case '?':
		s1 := s
		s1.first = '.'
		s2 := s
		s2.first = '#'
		// TODO: filter
		return []state{s1, s2}
	default:
		panic(s.first)
	}
}

func (s state) maybeAccumulate() (state, bool) {
	if s.acc > 0 {
		if s.expectPos >= s.expectLen {
			return s, false
		}

		// better match
		if s.acc != s.expect[s.expectPos] {
			return s, false
		}
		s.expectPos++
		s.acc = 0
	}
	return s, true
}

func newState(in string, expect []int) state {
	if len(in) == 0 {
		panic("here")
	}
	ret := state{
		first:     in[0],
		rest:      in[1:],
		acc:       0,
		expect:    [32]byte{},
		expectPos: 0,
		expectLen: len(expect),
	}
	for i, v := range expect {
		ret.expect[i] = byte(v)
	}
	return ret
}
