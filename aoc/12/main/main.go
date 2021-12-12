package main

import (
	"fmt"
	"strings"
)

var sample = `
start-A
start-b
A-c
A-b
b-d
A-end
b-end
`

var input = sample

type puzzle struct {
	edges map[string][]string
}

func (p puzzle) Walk() map[string]bool {
	ret := map[string]bool{}
	p.doWalk([]string{"start"}, false, ret)
	return ret
}

func (p puzzle) doWalk(path []string, hasDup bool, collect map[string]bool) {
	node := path[len(path)-1]
	for _, dst := range p.edges[node] {
		if dst == "start" {
			continue
		}
		newPath := append(path, dst)
		if dst == "end" {
			collect[strings.Join(newPath, ",")] = true
			continue
		}

		if dst == strings.ToUpper(dst) {
			// Big cave
			p.doWalk(newPath, hasDup, collect)
			continue
		}

		c := countTimes(dst, path)
		if c == 0 {
			p.doWalk(newPath, hasDup, collect)
			continue
		}

		if !hasDup && c == 1 {
			p.doWalk(newPath, true, collect)
			continue
		}

		// Don't visit
	}
}

func countTimes(needle string, haystack []string) int {
	count := 0
	for _, h := range haystack {
		if h == needle {
			count++
		}
	}
	return count
}

func main() {
	p := puzzle{
		edges: map[string][]string{},
	}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			panic(parts)
		}

		a, b := parts[0], parts[1]
		p.edges[a] = append(p.edges[a], b)
		p.edges[b] = append(p.edges[b], a)
	}
	fmt.Println(p.edges)

	results := p.Walk()
	for k := range results {
		fmt.Println(k)
	}
	fmt.Println(len(results))
}
