package main

import (
	"fmt"
	"strconv"
	"strings"
)

const sample = `2 3 0 3 10 11 12 1 1 0 1 99 2 1 1 2`

type node struct {
	child []*node
	meta  []int
}

func (n *node) metaSum() int {
	sum := 0
	for _, c := range n.child {
		sum += c.metaSum()
	}
	for _, m := range n.meta {
		sum += m
	}
	return sum
}

func (n *node) value() int {
	sum := 0
	if len(n.child) == 0 {
		for _, m := range n.meta {
			sum += m
		}
	} else {
		for _, m := range n.meta {
			if m < 1 || m > len(n.child) {
				continue // invalid
			}
			sum += n.child[m-1].value()
		}
	}

	return sum
}

func main() {
	tree := parse(sample)
	fmt.Println(tree.metaSum())
	fmt.Println(tree.value())
}

type input struct {
	s   []int
	pos int
}

func (in *input) parseNode() *node {
	nChildren := in.next()
	nMeta := in.next()
	if nMeta < 1 {
		panic(nMeta)
	}

	var n node
	for i := 0; i < nChildren; i++ {
		n.child = append(n.child, in.parseNode())
	}

	for i := 0; i < nMeta; i++ {
		n.meta = append(n.meta, in.next())
	}
	return &n
}

func (in *input) next() int {
	r := in.s[in.pos]
	in.pos++
	return r
}

func parse(text string) *node {
	parts := strings.Split(text, " ")
	in := input{
		s:   make([]int, len(parts)),
		pos: 0,
	}
	for i, p := range parts {
		in.s[i] = mustInt(p)
	}

	return in.parseNode()
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
