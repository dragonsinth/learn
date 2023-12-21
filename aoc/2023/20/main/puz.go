package main

import (
	"fmt"
	"os"
)

type puz struct {
	broadcast *node
	nodes     map[string]*node
}

func (p puz) initMemory() memory {
	return memory{vals: map[string]pulse{}}
}

func (p puz) press(mem memory, step int, debug bool) (int, int) {
	nLo, nHi := 0, 0
	queue := []message{{gen: 0, src: "button", dst: start, pulse: LO}}
	for i := 0; i < len(queue); i++ {
		m := queue[i]
		if debug {
			fmt.Println(m)
		}
		n := p.nodes[m.dst]
		p := m.pulse
		if p == LO {
			nLo++
		} else {
			nHi++
		}

		if m.dst == "rx" {
			if p == LO {
				os.Exit(1)
			}
		}

		switch n.wat {
		case PLAIN:
			// nothing special, just forward pulse
		case FF:
			if p == HI {
				continue // ignore
			}
			// whatever we sent last time, invert
			p = !mem.get(n.name, step)
		case CONJ:
			// count the high signals
			sum := 0
			for _, src := range n.src {
				if mem.get(src, step) == HI {
					sum++
				}
			}
			if sum == len(n.src) {
				p = LO
			} else {
				p = HI
			}

		default:
			panic(n.wat)
		}

		// record the last-sent pulse for this node
		mem.store(n.name, step, p)
		for _, dst := range n.dst {
			queue = append(queue, message{
				gen:   m.gen + 1,
				src:   n.name,
				dst:   dst,
				pulse: p,
			})
		}
	}
	return nLo, nHi
}

type memory struct {
	vals map[string]pulse
}

func (m memory) get(name string, step int) pulse {
	return m.vals[name]
}

func (m memory) store(name string, step int, p pulse) {
	was := m.vals[name]
	if was != p {
		m.vals[name] = p
	}
}

type node struct {
	name string
	wat  wat
	src  []string
	dst  []string
}

type wat byte

const (
	PLAIN = wat(0)
	FF    = wat('%')
	CONJ  = wat('&')
)

type pulse bool

const (
	LO = pulse(false)
	HI = pulse(true)
)

func (p pulse) String() string {
	if p == LO {
		return "low"
	} else {
		return "hi"
	}
}

type message struct {
	gen   int
	src   string
	dst   string
	pulse pulse
}

func (m message) String() string {
	return fmt.Sprintf("%d: %s -%s-> %s", m.gen, m.src, m.pulse, m.dst)
}
