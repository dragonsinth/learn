package main

import (
	"fmt"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"math/rand"
	"strings"
	"time"
)

func run(input string, iter int) {
	debug := iter <= 100
	rand.Seed(time.Now().UnixMicro())
	g := parseGraph(input)
	sz := len(g.nodes)
	histNode := make([]int, sz)
	histEdge := make([]int, len(g.edges))
	for i := 0; i < iter; i++ {
		src := rand.Intn(sz)
		dst := rand.Intn(sz)
		if src != dst {
			r := g.findShortestRoute(nodeId(src), nodeId(dst))
			if debug {
				fmt.Println(r)
			}
			if r != nil {
				for id, v := range r.seenNodes {
					if v && id != src && id != dst {
						histNode[id]++
					}
				}
				for id, v := range r.seenEdges {
					if v {
						histEdge[id]++
					}
				}
			}
		}
	}

	type result struct {
		name  string
		id    int
		score int
	}
	resultLess := func(a, b result) bool {
		return a.score > b.score
	}

	var results []result
	for id, v := range histNode {
		if v > 0 {
			results = append(results, result{
				name:  g.nodes[id].name,
				id:    id,
				score: v,
			})
		}
	}
	slices.SortFunc(results, resultLess)
	if debug {
		fmt.Println(results[:3])
	}

	results = nil
	for id, v := range histEdge {
		if v > 0 {
			results = append(results, result{
				name:  g.edges[id].name,
				id:    id,
				score: v,
			})
		}
	}
	slices.SortFunc(results, resultLess)
	if debug {
		fmt.Println(results[:3])
	}

	// remove the top edges
	for _, r := range results[:3] {
		e := g.edges[r.id]
		g.edgesByNodeId[e.src] = slices.DeleteFunc(g.edgesByNodeId[e.src], func(id edgeId) bool {
			return id == e.id
		})
		g.edgesByNodeId[e.dst] = slices.DeleteFunc(g.edgesByNodeId[e.dst], func(id edgeId) bool {
			return id == e.id
		})
	}

	// Count nodes reachable from 0
	reachable := g.countReachable(0)
	unreachable := sz - reachable
	fmt.Println(sz, reachable, unreachable, reachable*unreachable)
}

type graph struct {
	nodes []node
	edges []edge

	edgesByNodeId [][]edgeId
}

type node struct {
	name string
	id   nodeId
}

type edge struct {
	name string
	id   edgeId
	src  nodeId
	dst  nodeId
}

type nodeId int
type edgeId int

type route struct {
	loc       nodeId
	len       int
	seenNodes []bool
	seenEdges []bool
}

func (r route) String() string {
	var ids []int
	for i, v := range r.seenNodes {
		if v {
			ids = append(ids, i)
		}
	}
	return fmt.Sprintf("%d: %+v", r.len, ids)
}

func (g graph) findShortestRoute(start, end nodeId) *route {
	if start == end {
		return nil
	}

	seenNodes := map[nodeId]int{start: 0}
	startRoute := route{loc: start, len: 0, seenNodes: make([]bool, len(g.nodes)), seenEdges: make([]bool, len(g.edges))}
	startRoute.seenNodes[start] = true
	work := []route{startRoute}
	for len(work) > 0 {
		rt := work[0]
		work = work[1:]

		// find all the edges that start from node
		edgeIds := g.edgesByNodeId[rt.loc]
		for _, eId := range edgeIds {
			if rt.seenEdges[eId] {
				continue
			}

			e := g.edges[eId]
			var dst nodeId
			if e.src == rt.loc {
				dst = e.dst
			} else if e.dst == rt.loc {
				dst = e.src
			} else {
				panic("here")
			}
			if rt.seenNodes[dst] {
				continue
			}

			if v, ok := seenNodes[dst]; ok {
				if v > rt.len {
					panic("here")
				}
				continue // not shorted
			}
			seenNodes[dst] = rt.len

			newRt := rt
			newRt.loc = dst
			newRt.len++
			newRt.seenNodes = slices.Clone(rt.seenNodes)
			newRt.seenNodes[dst] = true
			newRt.seenEdges = slices.Clone(rt.seenEdges)
			newRt.seenEdges[eId] = true
			if dst == end {
				return &newRt // all done
			}
			work = append(work, newRt)
		}
	}
	return nil
}

func (g graph) countReachable(start nodeId) int {
	seenNodes := make([]bool, len(g.nodes))
	seenEdges := make([]bool, len(g.edges))
	work := []nodeId{start}
	for len(work) > 0 {
		src := work[0]
		work = work[1:]
		seenNodes[src] = true

		// find all the edges that start from node
		edgeIds := g.edgesByNodeId[src]
		for _, eId := range edgeIds {
			e := g.edges[eId]
			var dst nodeId
			if e.src == src {
				dst = e.dst
			} else if e.dst == src {
				dst = e.src
			} else {
				panic("here")
			}

			if seenNodes[dst] {
				continue
			}
			seenNodes[dst] = true
			seenEdges[eId] = true
			work = append(work, dst)
		}
	}
	sum := 0
	for _, v := range seenNodes {
		if v {
			sum++
		}
	}
	return sum
}

func parseGraph(input string) graph {
	nodeNames := map[string]bool{}
	var pairs []edgePair
	add := func(a, b string) {
		if a > b {
			a, b = b, a
		}
		pairs = append(pairs, edgePair{a, b})
		nodeNames[a] = true
		nodeNames[b] = true
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
	slices.SortFunc(pairs, func(p1, p2 edgePair) bool {
		if p1.a != p2.a {
			return p1.a < p2.a
		}
		return p1.b < p2.b
	})

	nodes := maps.Keys(nodeNames)
	slices.Sort(nodes)

	g := graph{
		nodes:         make([]node, len(nodes)),
		edges:         make([]edge, len(pairs)),
		edgesByNodeId: make([][]edgeId, len(nodes)),
	}
	nodeIds := map[string]nodeId{}
	for i, n := range nodes {
		g.nodes[i] = node{
			name: n,
			id:   nodeId(i),
		}
		nodeIds[n] = nodeId(i)
	}

	for i, p := range pairs {
		a, b := nodeIds[p.a], nodeIds[p.b]
		eId := edgeId(i)
		g.edges[i] = edge{
			name: p.a + p.b,
			id:   eId,
			src:  a,
			dst:  b,
		}
		g.edgesByNodeId[a] = append(g.edgesByNodeId[a], eId)
		g.edgesByNodeId[b] = append(g.edgesByNodeId[b], eId)
	}

	return g
}
