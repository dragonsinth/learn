package main

import (
	"bytes"
	"fmt"
	"github.com/dragonsinth/learn/aoc/termbox"
	"math"
	"os"
	"slices"
	"time"
)

type puz struct {
	units      map[pos]*unit
	walls      map[pos]bool
	w, h       int
	nElf, nGob int
	elfAtk     int
	turn       int

	term  termbox.Terminal
	debug bool
}

func (p *puz) Tick() bool {
	// run each unit in order
	units := mapValues(p.units)
	slices.SortFunc(units, func(a, b *unit) int {
		return readOrderPosLess(a.p, b.p)
	})

	for _, u := range units {
		if u.hp < 1 {
			// must have died
			continue
		}

		if p.nElf == 0 || p.nGob == 0 {
			return false
		}

		if pt := p.Move(*u); pt != nil {
			dst := *pt
			if !p.empty(dst) {
				panic(dst)
			}
			delete(p.units, u.p)
			u.p = dst
			p.units[u.p] = u
		}

		p.Attack(*u)
	}

	p.turn++
	return true
}

func (p *puz) Attack(u unit) {
	var targets []*unit
	for _, pt := range u.p.adjacent() {
		if foe := p.units[pt]; foe != nil && u.typ != foe.typ {
			targets = append(targets, foe)
		}
	}
	if len(targets) == 0 {
		return
	}
	slices.SortFunc(targets, func(a *unit, b *unit) int {
		if a.hp != b.hp {
			return a.hp - b.hp
		}
		if a.p.y != b.p.y {
			return a.p.y - b.p.y
		}
		return a.p.x - b.p.x
	})

	t := targets[0]
	switch u.typ {
	case ELF:
		t.hp -= p.elfAtk
	case GOBLIN:
		t.hp -= 3
	}
	if t.hp < 1 {
		delete(p.units, t.p)
		switch t.typ {
		case ELF:
			p.nElf--
		case GOBLIN:
			p.nGob--
		default:
			panic(t.typ)
		}
	}
}

func (p *puz) Move(u unit) *pos {
	// Find all target points.
	targets := make(map[pos]bool, p.nElf+p.nGob)
	for _, foe := range p.units {
		if u.typ == foe.typ {
			continue // friend
		}
		for _, pt := range foe.p.adjacent() {
			if pt == u.p {
				// there's an enemy right next to me
				return nil
			}
			if p.empty(pt) {
				targets[pt] = true
			}
		}
	}

	if len(targets) == 0 {
		return nil
	}

	type path []pos
	work := make([]path, 0, 1024)
	reachable := make(map[pos]path, len(targets))
	dist := math.MaxInt
	seen := make(map[pos]bool, p.w*p.h)
	for _, pt := range u.p.adjacent() {
		if p.empty(pt) {
			work = append(work, path{pt})
		}
	}

	for len(work) > 0 {
		pth := work[0]
		copy(work, work[1:])
		work = work[:len(work)-1]
		if len(pth) > dist {
			break
		}

		last := pth[len(pth)-1]
		if seen[last] {
			continue
		}
		seen[last] = true

		if targets[last] {
			// This is a valid target we haven't reached yet!
			reachable[last] = pth
			dist = len(pth)
			continue
		}

		// Consider any adjacent points.
		for _, pt := range last.adjacent() {
			if p.empty(pt) {
				newPath := make(path, 0, len(pth)+1)
				newPath = append(newPath, pth...)
				newPath = append(newPath, pt)
				work = append(work, newPath)
			}
		}
	}

	if len(reachable) == 0 {
		return nil
	}

	// Pick the best reachable target in read order
	reachPos := mapKeys(reachable)
	slices.SortFunc(reachPos, readOrderPosLess)
	pth := reachable[reachPos[0]]

	// walk that path
	return &pth[0]
}

func (p *puz) Print() {
	p.term.Render(p.Render(), os.Stdout)
}

func (p *puz) PrintFrame() {
	if p.debug {
		p.term.Render(p.Render(), os.Stdout)
		if p.term.Enabled() {
			p.wait()
		}
	}
}

func (p *puz) Render() [][]byte {
	var buf [][]byte
	for y := 0; y < p.h; y++ {
		var line []byte
		var rem bytes.Buffer
		for x := 0; x < p.w; x++ {
			var ch byte
			pt := pos{x, y}
			if u, ok := p.units[pt]; ok {
				ch = byte(u.typ)
				if rem.Len() == 0 {
					rem.WriteString("   ")
				} else {
					rem.WriteString(", ")
				}
				rem.WriteByte(ch)
				rem.WriteString(fmt.Sprintf("(%d)", u.hp))
			} else if p.walls[pt] {
				ch = '#'
			} else {
				ch = '.'
			}
			line = append(line, ch)
		}
		line = append(line, rem.Bytes()...)
		buf = append(buf, line)
	}
	buf = append(buf, []byte(fmt.Sprintf("Turn %d", p.turn)))
	return buf
}

func (p *puz) DebugRender(targets map[pos]bool, u *unit) [][]byte {
	var buf [][]byte
	for y := 0; y < p.h; y++ {
		var line []byte
		for x := 0; x < p.w; x++ {
			var ch byte
			pt := pos{x, y}
			if u != nil && u.p == pt {
				ch = byte(u.typ) - 'A' + 'a'
			} else if u, ok := p.units[pt]; ok {
				ch = byte(u.typ)
			} else if p.walls[pt] {
				ch = '#'
			} else if targets[pt] {
				ch = '?'
			} else {
				ch = '.'
			}
			line = append(line, ch)
		}
		buf = append(buf, line)
	}
	buf = append(buf, []byte(fmt.Sprintf("Turn %d", p.turn)))
	return buf
}

func (p *puz) empty(pt pos) bool {
	_, hasUnit := p.units[pt]
	hasWall := p.walls[pt]
	return !hasUnit && !hasWall
}

func (p *puz) wait() {
	time.Sleep(100 * time.Millisecond)
}

func (p *puz) hpSum() int {
	sum := 0
	for _, u := range p.units {
		sum += u.hp
	}
	return sum
}

func mapValues[K comparable, V any](in map[K]V) []V {
	var ret []V
	for _, v := range in {
		ret = append(ret, v)
	}
	return ret
}
