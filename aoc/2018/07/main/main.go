package main

import (
	"fmt"
	"golang.org/x/exp/slices"
	"regexp"
	"strings"

	_ "golang.org/x/exp/slices"
)

const sample = `
Step C must be finished before step A can begin.
Step C must be finished before step F can begin.
Step A must be finished before step B can begin.
Step A must be finished before step D can begin.
Step B must be finished before step E can begin.
Step D must be finished before step E can begin.
Step F must be finished before step E can begin.
`

const (
	nWorkers = 2
	addTime  = 0
	input    = sample
)

//const (
//	nWorkers = 5
//	addTime  = 60
//	input    = data
//)

var re = regexp.MustCompile(`^Step ([A-Z]) must be finished before step ([A-Z]) can begin.$`)

type id string

type node struct {
	id      id
	up      []*node
	dn      []*node
	started bool
	done    bool
}

func (n *node) Ready() bool {
	if n.started || n.done {
		return false
	}
	for _, up := range n.up {
		if !up.done {
			return false
		}
	}
	return true
}

type rule struct {
	up, dn id
}

func (r rule) String() string {
	return string(r.up) + "->" + string(r.dn)
}

type worker struct {
	n    *node
	when int
}

func (w worker) String() string {
	if w.n == nil {
		return "."
	}
	return string(w.n.id)
}

func (w worker) Format(f fmt.State, verb rune) {
	if w.n == nil {
		_, _ = f.Write([]byte("."))
	} else {
		_, _ = f.Write([]byte(w.n.id))
	}
}

func main() {
	nodes := parse()

	// part 1
	var sol string
	for len(sol) < len(nodes) {
		ready := readyNode(nodes)
		ready.started = true
		ready.done = true
		sol = sol + string(ready.id)
	}
	fmt.Println(sol)

	// part 2
	nodes = parse()

	sol = ""
	workers := make([]worker, nWorkers)
	for now := 0; len(sol) < len(nodes); now++ {
		// finish any done workers
		for i := range workers {
			w := &workers[i]
			if w.n != nil && w.when == now {
				sol = sol + string(w.n.id)
				w.n.done = true
				w.n = nil
			}
		}

		// each free worker should try to acquire a ready node
		for i := range workers {
			w := &workers[i]
			if w.n == nil {
				if n := readyNode(nodes); n != nil {
					n.started = true
					w.n = n
					w.when = now + addTime + (int(n.id[0]) - 'A' + 1)
				}
			}
		}

		fmt.Printf("%d %+v %s\n", now, workers, sol)
	}
}

func readyNode(nodes map[id]*node) *node {
	var ready []*node
	for _, n := range nodes {
		if n.Ready() {
			ready = append(ready, n)
		}
	}
	if len(ready) == 0 {
		return nil
	}

	slices.SortFunc(ready, func(a, b *node) bool {
		return a.id < b.id
	})
	return ready[0]
}

func parse() map[id]*node {
	nodes := map[id]*node{}
	var rules []rule
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !re.MatchString(line) {
			panic(line)
		}
		m := re.FindStringSubmatch(line)
		upId, dnId := m[1][0], m[2][0]
		r := rule{
			up: id(upId),
			dn: id(dnId),
		}
		rules = append(rules, r)
		nodes[(&node{id: r.up}).id] = &node{id: r.up}
		nodes[(&node{id: r.dn}).id] = &node{id: r.dn}
	}

	for _, r := range rules {
		up := nodes[r.up]
		dn := nodes[r.dn]
		up.dn = append(up.dn, dn)
		dn.up = append(dn.up, up)
	}

	return nodes
}
