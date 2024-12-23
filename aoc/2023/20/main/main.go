package main

import (
	"fmt"
	"strings"
)

var samples = []string{
	`
broadcaster -> a, b, c
%a -> b
%b -> c
%c -> inv
&inv -> a
`, `
broadcaster -> a
%a -> inv, con
&inv -> b
%b -> con
&con -> output
`,
}

const start = "broadcaster"

func main() {
	part1(samples[0], 1, true)
	part1(samples[1], 4, true)

	fmt.Println()

	part1(samples[0], 1000, false)
	part1(samples[1], 1000, false)

	//part1(data, 1000, false)
	//part2(data, 5000)
}

func part1(input string, count int, debug bool) {
	p := parse(input)
	mem := p.initMemory()
	sumLo, sumHi := 0, 0
	for step := 0; step < count; step++ {
		nLo, nHi := p.press(mem, step, debug)
		if debug {
			fmt.Println(nLo, nHi)
		}
		sumLo += nLo
		sumHi += nHi
	}
	fmt.Println(sumLo, sumHi, sumLo*sumHi)
}

func part2(input string, count int) {
	p := parse(input)
	mem := p.initMemory()
	for step := 0; step < count; step++ {
		md, ok := p.memDump(mem, step)
		if !ok {
			fmt.Println(step, md, ok)
			return
		}
		p.press(mem, step, false)
	}
	md, ok := p.memDump(mem, count)
	fmt.Println(count, md, ok)
}

func parse(input string) puz {
	p := puz{
		broadcast: nil,
		nodes:     map[string]*node{},
	}

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " -> ")
		if len(parts) != 2 {
			panic(parts)
		}
		name := parts[0]
		w := PLAIN
		if name[0] == '%' {
			w = FF
			name = name[1:]
		} else if name[0] == '&' {
			w = CONJ
			name = name[1:]
		}
		out := strings.Split(parts[1], ", ")
		n := &node{
			name: name,
			wat:  w,
			dst:  out,
		}
		p.nodes[name] = n
		if name == start {
			p.broadcast = n
		}
	}

	// go back and update the src nodes
	for _, n := range p.nodes {
		for _, dst := range n.dst {
			dn, ok := p.nodes[dst]
			if !ok {
				dn = &node{name: dst, wat: 0}
				p.nodes[dst] = dn
			}
			dn.src = append(dn.src, n.name)
		}
	}
	return p
}

func mapKeys[K comparable, V any](in map[K]V) []K {
	var ret []K
	for k := range in {
		ret = append(ret, k)
	}
	return ret
}
