package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/sliceheap"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var sample = `
Valve AA has flow rate=0; tunnels lead to valves DD, II, BB
Valve BB has flow rate=13; tunnels lead to valves CC, AA
Valve CC has flow rate=2; tunnels lead to valves DD, BB
Valve DD has flow rate=20; tunnels lead to valves CC, AA, EE
Valve EE has flow rate=3; tunnels lead to valves FF, DD
Valve FF has flow rate=0; tunnels lead to valves EE, GG
Valve GG has flow rate=0; tunnels lead to valves FF, HH
Valve HH has flow rate=22; tunnel leads to valve GG
Valve II has flow rate=0; tunnels lead to valves AA, JJ
Valve JJ has flow rate=21; tunnel leads to valve II
`

var (
	re = regexp.MustCompile(`^Valve ([A-Z][A-Z]) has flow rate=(\d+); tunnels? leads? to valves? (([A-Z][A-Z])(, [A-Z][A-Z])*)$`)
)

func main() {
	part1(parsePuzzle(sample), 1651)
	fmt.Println()
	part2(parsePuzzle(sample), 1707)
	fmt.Println()
}

func parsePuzzle(input string) puzzle {
	p := puzzle{nodesByName: map[string]*node{}}

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !re.MatchString(line) {
			panic(line)
		}

		matches := re.FindStringSubmatch(line)
		var dstNames []string
		for _, v := range strings.Split(matches[3], ", ") {
			dstNames = append(dstNames, v)
		}
		p.nodes = append(p.nodes, node{
			name:     matches[1],
			flow:     mustInt(matches[2]),
			dstNames: dstNames,
		})
	}

	// resolve
	for i := range p.nodes {
		n := &p.nodes[i]
		p.nodesByName[n.name] = n
		if n.flow > 0 {
			n.switchIdx = p.nSwitches
			p.nSwitches++
		}
	}
	for i := range p.nodes {
		n := &p.nodes[i]
		for _, dstName := range n.dstNames {
			dst := p.nodesByName[dstName]
			if dst == nil {
				panic(dstName)
			}
			n.dsts = append(n.dsts, dst)
		}
	}
	p.nodesByFlow = append([]node{}, p.nodes...)
	sort.Slice(p.nodesByFlow, func(i, j int) bool {
		return p.nodesByFlow[i].flow > p.nodesByFlow[j].flow
	})
	return p
}

func part1(p puzzle, want int) {
	start := p.nodesByName["AA"]
	j := &journey{
		player:   start,
		elephant: nil,
		turn:     0,
		maxTurns: 30,
		score:    0,
		step:     "start AA",
		switches: make([]bool, p.nSwitches),
	}
	j.maxScore = p.maxTheoretical(j)
	fmt.Println("max=", j.maxScore)

	best := p.computeBestJourney(j)
	for n := best; n != nil; n = n.prev {
		fmt.Println(n.step)
	}
	fmt.Println("best=", best.score)
	if want != 0 && want != best.score {
		log.Fatalf("want=%d, got=%d", want, best.score)
	}
}

func part2(p puzzle, want int) {
	start := p.nodesByName["AA"]
	j := &journey{
		player:   start,
		elephant: start,
		turn:     0,
		maxTurns: 26,
		score:    0,
		step:     "both start AA",
		switches: make([]bool, p.nSwitches),
	}
	j.maxScore = p.maxTheoreticalWithElephant(j)
	fmt.Println("max=", j.maxScore)

	best := p.computeBestJourneyWithElephant(j)
	for n := best; n != nil; n = n.prev {
		fmt.Println(n.step)
	}
	fmt.Println("best=", best.score)
	if want != 0 && want != best.score {
		log.Fatalf("want=%d, got=%d", want, best.score)
	}
}

type journey struct {
	prev *journey

	player       *node
	elephant     *node
	elephantTurn bool

	turn     int
	maxTurns int
	score    int
	maxScore int

	step      string
	switches  []bool
	nSwitches int
	switchKey string
}

func (j *journey) clone() *journey {
	cp := *j
	cp.switches = append([]bool{}, j.switches...)
	return &cp
}

func (j *journey) child() *journey {
	cp := *j
	cp.prev = j
	cp.switches = append([]bool{}, j.switches...)
	return &cp
}

func (j *journey) playerFlip(loc *node) *journey {
	if j.switches[loc.switchIdx] {
		panic(loc)
	}

	r := j.child()
	r.switches[loc.switchIdx] = true
	r.nSwitches++
	r.step = "player open " + loc.name
	r.score += (r.maxTurns - r.turn - 1) * loc.flow
	if r.score > r.maxScore {
		panic("here")
	}
	if j.elephant == nil {
		r.turn++
	} else {
		r.elephantTurn = true
	}

	// recompute switch key
	buf := make([]byte, len(r.switches))
	for i, v := range r.switches {
		if v {
			buf[i] = '+'
		} else {
			buf[i] = '-'
		}
	}
	r.switchKey = string(buf)
	return r
}

func (j *journey) elephantFlip(loc *node) *journey {
	if j.switches[loc.switchIdx] {
		panic(loc)
	}

	r := j.child()
	r.switches[loc.switchIdx] = true
	r.nSwitches++
	r.step = "elephant open " + loc.name
	r.score += (r.maxTurns - r.turn - 1) * loc.flow
	if r.score > r.maxScore {
		panic("here")
	}
	r.turn++
	r.elephantTurn = false

	// recompute switch key
	buf := make([]byte, len(r.switches))
	for i, v := range r.switches {
		if v {
			buf[i] = '+'
		} else {
			buf[i] = '-'
		}
	}
	r.switchKey = string(buf)
	return r
}

func (j *journey) playerMove(dst *node) *journey {
	if j.elephantTurn {
		panic("not player turn")
	}
	r := j.child()
	if j.elephant == nil {
		r.turn++
	} else {
		r.elephantTurn = true
	}
	r.player = dst
	r.step = "player moves " + dst.name
	return r
}

func (j *journey) elephantMove(dst *node) *journey {
	if !j.elephantTurn {
		panic("not elephant turn")
	}
	r := j.child()
	r.turn++
	r.elephantTurn = false
	r.elephant = dst
	r.step = "elephant moves " + dst.name
	return r
}

func (j *journey) key() journeyKey {
	ret := journeyKey{
		score:        j.score,
		turn:         j.turn,
		elephantTurn: j.elephantTurn,
		pLoc:         j.player.name,
		switchKey:    j.switchKey,
	}
	if j.elephant != nil {
		ret.eLoc = j.elephant.name
	}
	return ret
}

type journeyKey struct {
	score        int
	turn         int
	elephantTurn bool
	pLoc         string
	eLoc         string
	switchKey    string
}

func (p *puzzle) maxTheoretical(j *journey) int {
	if j.nSwitches == p.nSwitches || j.turn == j.maxTurns {
		return j.score
	}

	score := j.score
	ploc := j.player
	switches := append([]bool{}, j.switches...)
	nSwitches := j.nSwitches

	// If the player starts on a switch, we need to compute this two ways:
	// - one where they take the switch
	// - one where they don't
	var altScore int
	if ploc.flow > 0 && !j.switches[ploc.switchIdx] {
		// force the player to move by resetting them to start
		cp := j.clone()
		cp.player = p.nodesByName["AA"]
		altScore = p.maxTheoretical(cp)
	}

	turn := j.turn
	for {
		// If we're sitting on a live node, flip it.
		if ploc.flow > 0 && !switches[ploc.switchIdx] {
			switches[ploc.switchIdx] = true
			score += ploc.flow * (j.maxTurns - turn - 1)
			nSwitches++
			if nSwitches == p.nSwitches {
				return max(score, altScore)
			}
		} else {
			// Otherwise warp to the best node.
			best := p.bestToFlip(switches, nil)
			if best == nil {
				panic("here")
			}
			ploc = best
		}

		turn++
		if turn == j.maxTurns {
			return max(score, altScore)
		}
	}
}

func (p *puzzle) maxTheoreticalWithElephant(j *journey) int {
	if j.nSwitches == p.nSwitches || j.turn == j.maxTurns {
		return j.score
	}

	score := j.score
	ploc := j.player
	eloc := j.elephant
	switches := append([]bool{}, j.switches...)
	nSwitches := j.nSwitches

	// If the player and/or elephant starts on a lock, we need to compute this multiple ways.
	// - ones where they take the lock
	// - ones where they don't
	var altScore int
	playerBadSwitch := ploc.flow > 0 && !j.switches[ploc.switchIdx] && ploc != p.bestToFlip(switches, eloc)
	elephantBadSwitch := eloc.flow > 0 && !j.switches[eloc.switchIdx] && eloc != p.bestToFlip(switches, ploc)
	if playerBadSwitch {
		// force the player to move by resetting them to start
		cp := j.clone()
		cp.player = p.nodesByName["AA"]
		altScore = max(altScore, p.maxTheoreticalWithElephant(cp))
	}
	if elephantBadSwitch {
		// force the player to move by resetting them to start
		cp := j.clone()
		cp.elephant = p.nodesByName["AA"]
		altScore = max(altScore, p.maxTheoreticalWithElephant(cp))
	}
	if playerBadSwitch && elephantBadSwitch {
		// force the player to move by resetting them to start
		cp := j.clone()
		cp.player = p.nodesByName["AA"]
		cp.elephant = p.nodesByName["AA"]
		altScore = max(altScore, p.maxTheoreticalWithElephant(cp))
	}

	turn := j.turn
	elephantTurn := j.elephantTurn
	for {
		if !elephantTurn {
			// Player turn

			// If we're sitting on a live node, flip it.
			if ploc.flow > 0 && !switches[ploc.switchIdx] {
				switches[ploc.switchIdx] = true
				score += ploc.flow * (j.maxTurns - turn - 1)
				nSwitches++
				if nSwitches == p.nSwitches {
					return max(altScore, score)
				}
			} else {
				// Otherwise warp to the best node.
				best := p.bestToFlip(switches, eloc)
				if best != nil {
					ploc = best
				}
			}

			elephantTurn = true
		} else {
			// Elephant turn

			// If we're sitting on a live node, flip it.
			if eloc.flow > 0 && !switches[eloc.switchIdx] {
				switches[eloc.switchIdx] = true
				score += eloc.flow * (j.maxTurns - turn - 1)
				nSwitches++
				if nSwitches == p.nSwitches {
					return max(altScore, score)
				}
			} else {
				// Otherwise warp to the best node.
				best := p.bestToFlip(switches, ploc)
				if best != nil {
					eloc = best
				}
			}

			elephantTurn = false
			turn++
			if turn == j.maxTurns {
				return max(altScore, score)
			}
		}
	}
}

func (p *puzzle) bestToFlip(switches []bool, avoid *node) *node {
	var best *node
	for i := range p.nodesByFlow {
		n := &p.nodesByFlow[i]
		if n.flow == 0 {
			break
		}
		if switches[n.switchIdx] {
			continue // already flipped
		}
		if avoid != nil && n.name == avoid.name {
			continue // avoid this one
		}
		best = n
		break
	}
	return best
}

func betterJourney(a *journey, b *journey) bool {
	if a.maxScore != b.maxScore {
		return a.maxScore > b.maxScore
	}
	if a.score != b.score {
		return a.score > b.score
	}
	if a.turn != b.turn {
		return a.turn > b.turn
	}
	return false
}

func (p *puzzle) computeBestJourney(j *journey) *journey {
	seen := map[journeyKey]bool{}
	work := sliceheap.New(betterJourney)

	push := func(j *journey) {
		k := j.key()
		if !seen[k] {
			j.maxScore = p.maxTheoretical(j)
			if j.maxScore < j.score {
				fmt.Println(p.maxTheoretical(j))
				panic("here")
			}
			if j.prev != nil && j.maxScore > j.prev.maxScore {
				fmt.Println(p.maxTheoretical(j.prev))
				fmt.Println(p.maxTheoretical(j))
				panic("here")
			}
			work.Push(j)
			seen[k] = true
		}
	}

	push(j)
	trackMax := j.maxScore
	for {
		j := work.Pop()
		if j.turn == j.maxTurns || j.nSwitches == p.nSwitches {
			return j
		}

		if j.maxScore < trackMax {
			trackMax = j.maxScore
			fmt.Println(trackMax, len(seen), work.Len())
			seen = map[journeyKey]bool{}
		}

		// Consider flipping on the current switch.
		loc := j.player
		if loc.flow > 0 && !j.switches[loc.switchIdx] {
			// Add a work item that flips on the current switch.
			push(j.playerFlip(loc))
		}

		// Traverse each connection.
		for _, dst := range loc.dsts {
			push(j.playerMove(dst))
		}
	}
}

func (p *puzzle) computeBestJourneyWithElephant(j *journey) *journey {
	seen := map[journeyKey]bool{}
	work := sliceheap.New(betterJourney)

	push := func(j *journey) {
		k := j.key()
		if !seen[k] {
			j.maxScore = p.maxTheoreticalWithElephant(j)
			if j.maxScore < j.score {
				fmt.Println(p.maxTheoreticalWithElephant(j))
				panic("here")
			}
			if j.prev != nil && j.maxScore > j.prev.maxScore {
				fmt.Println(p.maxTheoreticalWithElephant(j.prev))
				fmt.Println(p.maxTheoreticalWithElephant(j))
				panic("here")
			}
			work.Push(j)
			seen[k] = true
		}
	}

	push(j)

	trackMax := j.maxScore
	for {
		j := work.Pop()
		if j.turn == j.maxTurns || j.nSwitches == p.nSwitches {
			return j
		}

		if j.maxScore < trackMax {
			trackMax = j.maxScore
			fmt.Println(trackMax, len(seen), work.Len())
			seen = map[journeyKey]bool{}
		}

		if !j.elephantTurn {
			// Consider flipping on the current switch.
			ploc := j.player
			if ploc.flow > 0 && !j.switches[ploc.switchIdx] {
				// Add a work item that flips on the current switch.
				push(j.playerFlip(ploc))
			}

			// Traverse each connection.
			for _, dst := range ploc.dsts {
				push(j.playerMove(dst))
			}
		} else {
			// Consider flipping on the current switch.
			eloc := j.elephant
			if eloc.flow > 0 && !j.switches[eloc.switchIdx] {
				// Add a work item that flips on the current switch.
				push(j.elephantFlip(eloc))
			}

			// Traverse each connection.
			for _, dst := range eloc.dsts {
				push(j.elephantMove(dst))
			}
		}
	}
}

type node struct {
	name      string
	flow      int
	switchIdx int
	dstNames  []string
	dsts      []*node
}

type puzzle struct {
	nodes       []node
	nodesByName map[string]*node
	nodesByFlow []node
	nSwitches   int
}

func max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
