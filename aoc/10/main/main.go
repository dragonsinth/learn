package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var sample = `
[({(<(())[]>[[{[]{<()<>>
[(()[<>])]({[<{<<[]>>(
{([(<{}[<>[]}>{[]{[(<()>
(((({<>}<{<{<>}{[]{[]{}
[[<[([]))<([[{}[[()]]]
[{[{({}]{}}([{[{{{}}([]
{<[[]]>}<{[{[{[]{()[[[]
[<(<(<(<{}))><([]([]()
<{([([[(<>()){}]>(<<{{
<{([{{}}[<[[[<>{}]]]>[]]
`

var input = sample

var (
	corruptPoints = map[byte]int{
		')': 3,
		']': 57,
		'}': 1197,
		'>': 25137,
	}

	incompletePoints = map[byte]int{
		')': 1,
		']': 2,
		'}': 3,
		'>': 4,
	}
)

type stream struct {
	b []byte
	p int
}

func main() {
	corruptScore := 0
	var incompleteScores []int
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		s := &stream{
			b: []byte(line),
			p: 0,
		}

		if err := s.parseChunks(); err != nil {
			if ce, ok := err.(errCorrupted); ok {
				pts := corruptPoints[ce.found]
				corruptScore += pts
				fmt.Println(line, err, pts)
			} else if ie, ok := err.(errIncomplete); ok {
				pts := 0
				for _, c := range ie.need {
					pts *= 5
					pts += incompletePoints[c]
				}
				fmt.Println(line, err, pts)
				incompleteScores = append(incompleteScores, pts)
			} else {
				panic(err)
			}
		} else {
			fmt.Println(line, "pass")
		}
	}
	fmt.Println(corruptScore)
	sort.Ints(incompleteScores)
	fmt.Println(incompleteScores[len(incompleteScores)/2])
}

var (
	errUnexpected = errors.New("unexpected")

	match = map[byte]byte{
		'(': ')',
		'[': ']',
		'{': '}',
		'<': '>',
	}
)

type errCorrupted struct {
	expect byte
	found  byte
}

func (e errCorrupted) Error() string {
	return fmt.Sprintf("Corrupt: expected %c, but found %c instead.", e.expect, e.found)
}

type errIncomplete struct {
	need []byte
}

func (e errIncomplete) Error() string {
	return fmt.Sprintf("Incomplete: need %s", string(e.need))
}

func (s *stream) parseChunk() error {
	if s.p == len(s.b) {
		return nil
	}

	ch1 := s.b[s.p]
	switch ch1 {
	case '(', '[', '{', '<':
		s.p++
	case ')', ']', '}', '>':
		return nil // hand this back to the caller
	default:
		return errUnexpected
	}

	need := match[ch1]
	if err := s.parseChunks(); err != nil {
		if ie, ok := err.(errIncomplete); ok {
			ie.need = append(ie.need, need)
			return ie
		}
		return err
	} else if s.p == len(s.b) {
		return errIncomplete{need: []byte{need}}
	}

	ch2 := s.b[s.p]
	switch ch2 {
	case need:
		s.p++
		return nil
	case ')', ']', '}', '>':
		return errCorrupted{
			expect: need,
			found:  ch2,
		}
	default:
		return errUnexpected
	}
}

func (s *stream) parseChunks() error {
	for {
		if s.p == len(s.b) {
			return nil
		}

		ch := s.b[s.p]
		switch ch {
		case '(', '[', '{', '<':
			if err := s.parseChunk(); err != nil {
				return err
			}
		case ')', ']', '}', '>':
			return nil // hand this back to the caller
		default:
			return errUnexpected
		}
	}
}
