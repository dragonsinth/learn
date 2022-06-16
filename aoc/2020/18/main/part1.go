package main

import "fmt"

func parse1(line string) int {
	s := &stream{line: line, pos: 0}
	val := parse1Expr(s)
	if !s.Eof() {
		panic(fmt.Sprintf("unconsumed: %q", s.line[s.pos:]))
	}
	return val
}

func parse1Expr(s *stream) int {
	val := parse1Terminal(s)
	for {
		if s.Eof() || s.Peek() == ')' {
			return val
		}

		s.MustPop(' ')
		c := s.Pop()
		switch c {
		case '*':
		case '+':
		default:
			panic(fmt.Sprintf("expected op, got %s", string(c)))
		}
		s.MustPop(' ')
		rhs := parse1Terminal(s)
		switch c {
		case '*':
			val = val * rhs
		case '+':
			val = val + rhs
		default:
			panic("bug")
		}
	}
}

func parse1Terminal(s *stream) int {
	if s.Peek() == '(' {
		return parse1Paren(s)
	} else {
		return parse1Digit(s)
	}
}

func parse1Digit(s *stream) int {
	c := s.Pop()
	if c < '0' || c > '9' {
		panic(fmt.Sprintf("expected digit, got %s", string(c)))
	}
	return int(c - '0')
}

func parse1Paren(s *stream) int {
	s.MustPop('(')
	val := parse1Expr(s)
	s.MustPop(')')
	return val
}
