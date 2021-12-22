package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `
[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]
[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]
[[2,[[0,8],[3,4]]],[[[6,7],1],[7,[1,6]]]]
[[[[2,4],7],[6,[0,5]]],[[[6,8],[2,8]],[[2,1],[4,5]]]]
[7,[5,[[3,8],[1,4]]]]
[[2,[2,2]],[8,[8,1]]]
[2,9]
[1,[[[9,3],9],[[9,0],[0,7]]]]
[[[5,[7,4]],7],1]
[[[[4,2],2],6],[8,7]]
`

var sample2 = `
[[[0,[5,8]],[[1,7],[9,6]]],[[4,[1,2]],[[1,4],2]]]
[[[5,[2,8]],4],[5,[[9,9],0]]]
[6,[[[6,2],[5,6]],[[7,6],[4,7]]]]
[[[6,[0,7]],[0,9]],[4,[9,[9,0]]]]
[[[7,[6,4]],[3,[1,3]]],[[[5,5],1],9]]
[[6,[[7,3],[3,2]]],[[[3,8],[5,7]],4]]
[[[[5,4],[7,7]],8],[[8,3],8]]
[[9,3],[[9,9],[6,[4,9]]]]
[[2,[[7,7],7]],[[5,8],[[9,3],[0,2]]]]
[[[[5,2],5],[8,[3,7]]],[[5,[7,5]],[4,4]]]
`

var input = sample2

func main() {
	var nums []*Number
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		n := parse(line)
		nums = append(nums, n)
	}

	best := 0
	for i := range nums {
		for j := range nums {
			if i == j {
				continue
			}
			m := reduce(pair(nums[i], nums[j])).Magnitude()
			if m > best {
				best = m
			}
		}
	}
	fmt.Println(best)
}

func main1() {
	var last *Number
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		current := parse(line)
		if last != nil {
			fmt.Println(" ", last)
			fmt.Println("+", current)
			current = reduce(pair(last, current))
			fmt.Println("=", current)
		}
		last = current
	}

	fmt.Println(last.Magnitude())
}

type Number struct {
	isLiteral bool
	literal   int

	left  *Number
	right *Number
}

func literal(v int) *Number {
	return &Number{
		isLiteral: true,
		literal:   v,
	}
}

func pair(left *Number, right *Number) *Number {
	return &Number{
		left:  left,
		right: right,
	}
}

func (p *Number) String() string {
	if p.isLiteral {
		return strconv.Itoa(p.literal)
	}
	return "[" + p.left.String() + "," + p.right.String() + "]"
}

func (p *Number) Magnitude() int {
	if p.isLiteral {
		return p.literal
	}
	return 3*p.left.Magnitude() + 2*p.right.Magnitude()
}

func explode(n *Number) *Number {
	ret, _, _ := explodeDepth(n, 0)
	return ret
}

func explodeDepth(n *Number, depth int) (*Number, int, int) {
	if n.isLiteral {
		return n, 0, 0
	}
	if depth == 4 {
		return literal(0), n.left.literal, n.right.literal
	}

	{
		l, exLeft, exRight := explodeDepth(n.left, depth+1)
		if l != n.left {
			// Replace n.left, bubble exLeft, add in exRight
			r := addLeft(n.right, exRight)
			return pair(l, r), exLeft, 0
		}
	}
	{
		r, exLeft, exRight := explodeDepth(n.right, depth+1)
		if r != n.right {
			// Replace n.right, bubble exRight, add in exLeft
			l := addRight(n.left, exLeft)
			return pair(l, r), 0, exRight
		}
	}

	return n, 0, 0
}

func addLeft(n *Number, v int) *Number {
	if n.isLiteral {
		return literal(n.literal + v)
	} else {
		return pair(addLeft(n.left, v), n.right)
	}
}

func addRight(n *Number, v int) *Number {
	if n.isLiteral {
		return literal(n.literal + v)
	} else {
		return pair(n.left, addRight(n.right, v))
	}
}

func split(n *Number) *Number {
	if n.isLiteral {
		if n.literal < 10 {
			return n
		}
		return pair(
			literal(n.literal/2),
			literal(n.literal/2+(n.literal%2)),
		)
	} else {
		if l := split(n.left); l != n.left {
			return pair(l, n.right)
		}
		if r := split(n.right); r != n.right {
			return pair(n.left, r)
		}
		return n
	}
}

func reduce(n *Number) *Number {
	for {
		newN := explode(n)
		if newN != n {
			n = newN
			continue
		}

		newN = split(n)
		if newN != n {
			n = newN
			continue
		}

		return n
	}
}

func parse(input string) *Number {
	s := &stream{
		b: []byte(input),
		p: 0,
	}

	ret := s.parseNumber()
	if s.p < len(s.b) {
		panic(fmt.Sprintf("unparsed input: %s", string(s.b[s.p:])))
	}
	return ret
}

type stream struct {
	b []byte
	p int
}

func (s *stream) Next() byte {
	r := s.b[s.p]
	s.p++
	return r
}

func (s *stream) parseNumber() *Number {
	if s.p == len(s.b) {
		return nil
	}

	if c := s.Next(); c == '[' {
		l := s.parseNumber()
		if c := s.Next(); c != ',' {
			panic(c)
		}
		r := s.parseNumber()
		if c := s.Next(); c != ']' {
			panic(c)
		}

		return &Number{
			left:  l,
			right: r,
		}
	} else if c >= '0' && c <= '9' {
		return literal(int(c - '0'))
	} else {
		panic(string(c))
	}
}
