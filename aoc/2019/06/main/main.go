package main

import (
	"fmt"
	"strings"
)

var sample = `
COM)B
B)C
C)D
D)E
E)F
B)G
G)H
D)I
E)J
J)K
K)L
K)YOU
I)SAN
`

var input = sample

type node struct {
	name   string
	target *node
	down   int
}

func main() {
	rules := map[string]string{}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ")")
		rules[parts[1]] = parts[0]
	}

	all := map[string]*node{
		"COM": {
			name:   "COM",
			target: nil,
			down:   0,
		},
	}

	for k := range rules {
		resolve(k, all, rules)
	}

	sum := 0
	for k, v := range all {
		fmt.Println(k, v.down)
		sum += v.down
	}
	fmt.Println(sum)

	fmt.Println(navigate(all["YOU"], all["SAN"]))
}

func navigate(n1 *node, n2 *node) int {
	if n1.target == n2.target {
		return 0
	}
	if n1.down > n2.down {
		return 1 + navigate(n1.target, n2)
	} else if n1.down < n2.down {
		return 1 + navigate(n1, n2.target)
	} else {
		return 2 + navigate(n1.target, n2.target)
	}
}

func resolve(k string, all map[string]*node, rules map[string]string) *node {
	if v, ok := all[k]; ok {
		return v
	}

	tgt := rules[k]
	if tgt == "" {
		panic(k)
	}

	tgtNode := resolve(tgt, all, rules)
	n := &node{
		name:   k,
		target: tgtNode,
		down:   tgtNode.down + 1,
	}
	all[k] = n
	return n
}
