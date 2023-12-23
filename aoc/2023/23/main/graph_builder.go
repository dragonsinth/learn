package main

import (
	"fmt"

	"golang.org/x/exp/slices"
)

func (p puz) computeGraph(start, dest pos, debug bool) graph {
	type segment struct {
		p1, p2 pos
		len    int
	}

	work := []task{{
		x:     start.x,
		y:     start.y + 1,
		d:     S,
		route: []taskKey{{start.x, start.y, S}},
	}}

	type arrow int
	const UN, FO, RV = arrow(0), arrow(1), arrow(2)

	canWalk := func(t task) (task, bool) {
		switch p.data[t.y][t.x] {
		case '#':
			return t, false // wall
		case '.', '^', '>', 'v', '<':
			return t, true // open
		default:
			panic(string(p.data[t.y][t.x]))
		}
	}

	checkArrow := func(t task) arrow {
		switch c := p.data[t.y][t.x]; c {
		case '.':
			return UN // open
		case '^', '>', 'v', '<':
			if t.d.sym() == c {
				return FO
			} else {
				return RV
			}
		default:
			panic(string(p.data[t.y][t.x]))
		}
	}

	nodes := map[pos]bool{
		start: true,
		dest:  true,
	}

	segs := map[segment]bool{}

	finishSegment := func(t task, ar arrow) {
		st := t.route[0]
		seg := segment{
			p1:  pos{st.x, st.y},
			p2:  pos{t.x, t.y},
			len: len(t.route),
		}
		switch ar {
		case FO:
		case RV:
			seg.p1, seg.p2 = seg.p2, seg.p1
		default:
			panic(ar)
		}

		if debug {
			fmt.Printf("found segment (%d,%d)->(%d,%d): %d (%t)\n", seg.p1.x, seg.p1.y, seg.p2.x, seg.p2.y, seg.len, segs[seg])
		}
		segs[seg] = true
	}

	for len(work) > 0 {
		cur := work[0]
		work = work[1:]

		// walk to the next intersection
		ar := UN
		for {
			if nodes[cur.pos()] {
				// found the end of a segment at a known intersection
				finishSegment(cur, ar)
				break
			}

			ar |= checkArrow(cur)

			next := make([]task, 0, 3)
			if n, ok := canWalk(cur.left()); ok {
				next = append(next, n)
			}
			if n, ok := canWalk(cur.right()); ok {
				next = append(next, n)
			}
			if n, ok := canWalk(cur.next()); ok {
				next = append(next, n)
			}

			if len(next) == 0 {
				panic("here")
			} else if len(next) == 1 {
				cur = next[0]
			} else {
				// Found an intersection, start new segments.
				nodes[cur.pos()] = true

				for _, n := range next {
					n.route = []taskKey{{cur.x, cur.y, n.d}}
					work = append(work, n)
				}

				finishSegment(cur, ar)
				break
			}
		}
	}

	g := graph{
		nodes:            nil,
		edges:            nil,
		edgesByNodeId:    make([][]edgeId, len(nodes)),
		edgesByNodeIdSrc: make([][]edgeId, len(nodes)),
		edgesByNodeIdDst: make([][]edgeId, len(nodes)),
	}

	for p := range nodes {
		g.nodes = append(g.nodes, node{
			id: 0,
			x:  p.x,
			y:  p.y,
		})
	}
	slices.SortFunc(g.nodes, func(a, b node) bool {
		if a.y != b.y {
			return a.y < b.y // ensures start at node 0, end at last node.
		}
		return a.x < b.x // doesn't matter
	})
	nodesByPos := map[pos]nodeId{}
	for i := range g.nodes {
		n := &g.nodes[i]
		n.id = nodeId(i)
		nodesByPos[n.pos()] = n.id
	}

	for s := range segs {
		g.edges = append(g.edges, edge{
			len: s.len,
			src: nodesByPos[s.p1],
			dst: nodesByPos[s.p2],
		})
	}
	for i := range g.edges {
		e := &g.edges[i]
		e.id = edgeId(i)
		g.edgesByNodeId[e.src] = append(g.edgesByNodeId[e.src], e.id)
		g.edgesByNodeId[e.dst] = append(g.edgesByNodeId[e.dst], e.id)
		g.edgesByNodeIdSrc[e.src] = append(g.edgesByNodeIdSrc[e.src], e.id)
		g.edgesByNodeIdDst[e.dst] = append(g.edgesByNodeIdDst[e.dst], e.id)
	}

	if debug {
		fmt.Println(len(g.nodes))
		fmt.Println(len(g.edges))
	}
	return g
}
