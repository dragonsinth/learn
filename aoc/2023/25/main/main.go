package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

var sample = `
jqt: rhn xhk nvd
rsh: frs pzl lsr
xhk: hfx
cmg: qnr nvd lhk bvb
rhn: xhk bvb hfx
bvb: xhk hfx
pzl: lsr hfx nvd
qnr: nvd
ntq: jqt hfx bvb xhk
nvd: lhk
lsr: lhk
rzs: qnr cmg lsr rsh
frs: qnr lhk lsr
`

func main() {
	visualize(sample, os.Stdout)
	run(sample, 100)
}

func visualize(input string, out io.Writer) {
	pairs, _ := parse(input)
	_, _ = fmt.Fprintln(out, `graph G {`)
	for _, p := range pairs {
		_, _ = fmt.Fprintf(out, `  %s -- %s [label="%s%s"];`, p.a, p.b, p.a, p.b)
		_, _ = fmt.Fprintln(out)
	}
	_, _ = fmt.Fprintln(out, `}`)
}

type edgePair struct {
	a, b string
}

func edgeSort(p1, p2 edgePair) int {
	if p1.a != p2.a {
		return strings.Compare(p1.a, p2.a)
	}
	return strings.Compare(p1.b, p2.b)
}

func parse(input string) ([]edgePair, []string) {
	nodes := map[string]bool{}
	var pairs []edgePair
	add := func(a, b string) {
		if b > a {
			a, b = b, a
		}
		pairs = append(pairs, edgePair{a, b})
		nodes[a] = true
		nodes[b] = true
	}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		src := parts[0]
		rest := strings.Fields(parts[1])
		for _, r := range rest {
			add(src, r)
		}
	}
	slices.SortFunc(pairs, edgeSort)

	n := mapKeys(nodes)
	slices.Sort(n)

	return pairs, n
}

func mapKeys[K comparable, V any](in map[K]V) []K {
	var ret []K
	for k := range in {
		ret = append(ret, k)
	}
	return ret
}
