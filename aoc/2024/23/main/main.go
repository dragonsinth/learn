package main

import (
	"fmt"
	"slices"
	"strings"
)

const sample = `
kh-tc
qp-kh
de-cg
ka-co
yn-aq
qp-ub
cg-tb
vc-aq
tb-ka
wh-tc
yn-cg
kh-ub
ta-co
de-co
tc-td
tb-wq
wh-td
ta-ka
td-qp
aq-cg
wq-ub
ub-vc
de-ta
wq-aq
wq-vc
wh-yn
ka-de
kh-ta
co-tc
wh-qp
tb-vc
td-yn
`

func main() {
	run(sample)
}

func run(input string) {
	conns := map[node][]node{}
	pairs := map[set]bool{}

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			panic(line)
		}
		a := newNode(parts[0])
		b := newNode(parts[1])
		if a > b {
			a, b = b, a
		}
		pairs[set{a, b}] = true
		conns[a] = append(conns[a], b)
		conns[b] = append(conns[b], a)
	}

	// sort
	for k, v := range conns {
		slices.Sort(v)
		conns[k] = v
	}

	// for every pair, see what the common conns are
	trios := increaseSets(conns, pairs, 2)

	sum := 0
	for r := range trios {
		if r[0].Hi() == 't' || r[1].Hi() == 't' || r[2].Hi() == 't' {
			fmt.Println(r[0], r[1], r[2])
			sum++
		}
	}
	fmt.Println(len(trios), sum)

	// see how deep the rabbit hole goes
	setSize := 3
	curSet := trios
	for len(curSet) > 1 {
		nextSet := increaseSets(conns, curSet, setSize)
		curSet = nextSet
		setSize++
	}
	for s := range curSet {
		var sb strings.Builder
		for i, e := range s[:setSize] {
			if i > 0 {
				sb.WriteRune(',')
			}
			sb.WriteString(e.String())
		}
		fmt.Println(setSize, sb.String())
	}
}

type set [64]node

func increaseSets(conns map[node][]node, in map[set]bool, inSize int) map[set]bool {
	outSize := inSize + 1
	// for incoming set, see what the common conns are
	ret := map[set]bool{}
	for e := range in {
		common := conns[e[0]]
		for i := 1; i < inSize; i++ {
			common = zipperMerge(common, conns[e[i]])
		}
		for _, c := range common {
			nextSet := e
			nextSet[inSize] = c
			slices.Sort(nextSet[:outSize])
			ret[nextSet] = true
		}
	}
	return ret
}

func zipperMerge(a []node, b []node) []node {
	var r []node
	ai, bi := 0, 0
	for ai < len(a) && bi < len(b) {
		if a[ai] < b[bi] {
			ai++
		} else if a[ai] > b[bi] {
			bi++
		} else if a[ai] == b[bi] {
			r = append(r, a[ai])
			ai++
			bi++
		} else {
			panic("here")
		}
	}
	return r
}

func newNode(in string) node {
	if len(in) != 2 {
		panic(in)
	}
	return node(in[0])<<8 | node(in[1])
}

type node uint16

func (n node) String() string {
	return string([]byte{n.Hi(), n.Lo()})
}

func (n node) Hi() byte {
	return byte((n >> 8) & 0xff)
}

func (n node) Lo() byte {
	return byte(n)
}
