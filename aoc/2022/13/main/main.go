package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var sample = `
[1,1,3,1,1]
[1,1,5,1,1]

[[1],[2,3,4]]
[[1],4]

[9]
[[8,7,6]]

[[4,4],4,4]
[[4,4],4,4,4]

[7,7,7,7]
[7,7,7]

[]
[3]

[[[]]]
[[]]

[1,[2,[3,[4,[5,6,7]]]],8,9]
[1,[2,[3,[4,[5,6,0]]]],8,9]
`

func main() {
	part1()
	part2()
}

func part1() {
	var last string
	index := 0
	sum := 0
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if last == "" {
			last = line
			continue
		}

		index++
		l := parse(last)
		last = ""
		r := parse(line)

		inOrder := compare(l, r)
		if inOrder < 1 {
			fmt.Printf("pair %d is inOrder\n", index)
			sum += index
		}
	}
	fmt.Println(sum)
}

func part2() {
	m1 := []any{[]any{2}}
	m2 := []any{[]any{6}}
	list := []any{
		m1,
		m2,
	}
	for _, line := range strings.Split(sample, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		list = append(list, parse(line))
	}
	sort.Slice(list, func(i, j int) bool {
		a := list[i]
		b := list[j]
		return compare(a, b) < 0
	})

	var i1, i2 int
	for i, e := range list {
		fmt.Println(e)
		if compare(e, m1) == 0 {
			i1 = i + 1
		}
		if compare(e, m2) == 0 {
			i2 = i + 1
		}
	}
	fmt.Println(i1, i2, i1*i2)
}

func compare(l any, r any) int {
	left, leftInt := l.(int)
	right, rightInt := r.(int)
	switch {
	case leftInt && rightInt:
		return left - right
	case leftInt && !rightInt:
		return compareList([]any{left}, r.([]any))
	case !leftInt && rightInt:
		return compareList(l.([]any), []any{right})
	case !leftInt && !rightInt:
		return compareList(l.([]any), r.([]any))
	default:
		panic(nil)
	}
}

func compareList(l []any, r []any) int {
	for {
		if len(l) == 0 && len(r) == 0 {
			return 0
		}
		if len(l) == 0 {
			return -1
		}
		if len(r) == 0 {
			return 1
		}

		c := compare(l[0], r[0])
		if c != 0 {
			return c
		}
		l = l[1:]
		r = r[1:]
	}
}

func parse(in string) any {
	p := &parser{
		in:  []byte(in),
		pos: 0,
	}
	return p.parseValue()
}

type parser struct {
	in  []byte
	pos int
}

func (p *parser) parseValue() any {
	if p.consume('[') {
		var ret []any
		if p.consume(']') {
			return ret
		}
		for {
			e := p.parseValue()
			ret = append(ret, e)
			if p.consume(',') {
				continue
			} else if p.consume(']') {
				return ret
			} else {
				panic(p.rem())
			}
		}
	}

	// must be a literal
	var buf []byte
	for p.peek() >= '0' && p.peek() <= '9' {
		buf = append(buf, p.pop())
	}
	return mustInt(string(buf))
}

func (p *parser) consume(c byte) bool {
	if p.in[p.pos] == c {
		p.pos++
		return true
	}
	return false
}

func (p *parser) rem() string {
	return string(p.in[p.pos:])
}

func (p *parser) peek() byte {
	return p.in[p.pos]
}

func (p *parser) pop() byte {
	ret := p.in[p.pos]
	p.pos++
	return ret
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
