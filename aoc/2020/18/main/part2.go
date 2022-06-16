package main

import "fmt"

func parse(line string) int {
	s := &stream{line: line, pos: 0}
	val := parseProduct(s)
	if !s.Eof() {
		panic(fmt.Sprintf("unconsumed: %q", s.line[s.pos:]))
	}
	return val
}

func parseProduct(s *stream) int {
	val := parseSum(s)
	for {
		if s.Eof() || s.Peek() == ')' {
			return val
		}

		s.MustPop('*')
		s.MustPop(' ')
		rhs := parseSum(s)
		val = val * rhs
	}
}

func parseSum(s *stream) int {
	val := parseTerminal(s)
	for {
		if s.Eof() || s.Peek() == ')' {
			return val
		}

		s.MustPop(' ')

		// If the next operator is a '+', keep going; if '*', we're done.
		c := s.Peek()
		switch c {
		case '*':
			// End of sum
			return val
		case '+':
			s.Pop()
			s.MustPop(' ')
			rhs := parseTerminal(s)
			val = val + rhs
		default:
			panic(fmt.Sprintf("expected op, got %s", string(c)))
		}
	}
}

func parseTerminal(s *stream) int {
	if s.Peek() == '(' {
		return parseParen(s)
	} else {
		return parseDigit(s)
	}
}

func parseDigit(s *stream) int {
	c := s.Pop()
	if c < '0' || c > '9' {
		panic(fmt.Sprintf("expected digit, got %s", string(c)))
	}
	return int(c - '0')
}

func parseParen(s *stream) int {
	s.MustPop('(')
	val := parseProduct(s)
	s.MustPop(')')
	return val
}
