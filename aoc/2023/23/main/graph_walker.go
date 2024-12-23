package main

import (
	"fmt"
	"slices"
)

type graph struct {
	nodes []node
	edges []edge

	edgesByNodeId    [][]edgeId
	edgesByNodeIdSrc [][]edgeId
	edgesByNodeIdDst [][]edgeId
}

type nodeId int

type node struct {
	id   nodeId
	x, y int
}

func (n node) pos() pos {
	return pos{n.x, n.y}
}

type edgeId int

type edge struct {
	id  edgeId
	len int
	src nodeId
	dst nodeId
}

type route struct {
	atNode    nodeId
	len       int
	seenNodes [40]bool
	seenEdges [64]bool
}

func (g graph) findAllRoutes(either bool, debug bool) []route {
	finalDst := nodeId(len(g.nodes) - 1)

	var complete []route
	work := []route{{
		atNode:    0, // start
		len:       0,
		seenNodes: [40]bool{0: true},
		seenEdges: [64]bool{},
	}}

	maybeTraverseEdge := func(rt route, e edge, dst nodeId) {
		if rt.seenNodes[dst] {
			return // already seen destination node
		}

		// walk this edge
		newT := route{
			atNode:    dst,
			len:       rt.len + e.len,
			seenNodes: rt.seenNodes,
			seenEdges: rt.seenEdges,
		}
		newT.seenEdges[e.id] = true
		newT.seenNodes[dst] = true

		if e.dst == finalDst {
			complete = append(complete, newT)
		} else {
			work = append(work, newT)
		}
	}

	edgeMap := g.edgesByNodeIdSrc
	if either {
		edgeMap = g.edgesByNodeId
	}

	for len(work) > 0 {
		rt := work[0]
		work = work[1:]

		// find all the edges that start from node
		edgeIds := edgeMap[rt.atNode]
		if debug {
			fmt.Printf("%d -> %+v\n", rt.atNode, edgeIds)
		}
		for _, eId := range edgeIds {
			if rt.seenEdges[eId] {
				continue // already seen this edge
			}
			e := g.edges[eId]
			// ok to try both; one will be ignored
			maybeTraverseEdge(rt, e, e.src)
			maybeTraverseEdge(rt, e, e.dst)
		}
	}
	slices.SortFunc(complete, func(a, b route) int {
		return b.len - a.len
	})
	return complete
}
