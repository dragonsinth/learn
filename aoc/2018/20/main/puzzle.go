package main

import "bytes"

type puzzle struct {
	wat      map[pos]room
	min, max pos
}

func (p *puzzle) Walk(in *input, loc pos) {
outer:
	for {
		p.min.x = min(p.min.x, loc.x)
		p.max.x = max(p.max.x, loc.x)
		p.min.y = min(p.min.y, loc.y)
		p.max.y = max(p.max.y, loc.y)

		switch c := in.Peek(); c {
		case 'N':
			in.Next()
			p.wat[loc] = p.wat[loc].With(N)
			loc = loc.Walk(N)
			p.wat[loc] = p.wat[loc].With(S)
		case 'E':
			in.Next()
			p.wat[loc] = p.wat[loc].With(E)
			loc = loc.Walk(E)
			p.wat[loc] = p.wat[loc].With(W)
		case 'S':
			in.Next()
			p.wat[loc] = p.wat[loc].With(S)
			loc = loc.Walk(S)
			p.wat[loc] = p.wat[loc].With(N)
		case 'W':
			in.Next()
			p.wat[loc] = p.wat[loc].With(W)
			loc = loc.Walk(W)
			p.wat[loc] = p.wat[loc].With(E)
		case '(':
			in.Next()
			for {
				// recurse
				p.Walk(in, loc)

				switch c := in.Peek(); c {
				case ')':
					in.Next()
					continue outer
				case '|':
					in.Next()
					continue
				default:
					panic(string(c))
				}
			}
		case ')':
			return
		case '|':
			return
		case '$':
			return
		default:
			panic(string(c))
		}
	}
}

func (p *puzzle) Render() [][]byte {
	w := p.max.x - p.min.x + 1
	var ret [][]byte
	for y := p.min.y; y <= p.max.y; y++ {
		// north wall
		buf := make([]byte, 0, 2*w+1)
		for x := p.min.x; x <= p.max.x; x++ {
			buf = append(buf, '#')
			if p.wat[pos{x, y}].Has(N) {
				buf = append(buf, '-')
			} else {
				buf = append(buf, '#')
			}
		}
		// east wall
		buf = append(buf, '#')
		ret = append(ret, buf)

		// room center, west door
		buf = make([]byte, 0, 2*w+1)
		for x := p.min.x; x <= p.max.x; x++ {
			if p.wat[pos{x, y}].Has(W) {
				buf = append(buf, '|')
			} else {
				buf = append(buf, '#')
			}
			if x == 0 && y == 0 {
				buf = append(buf, 'X')
			} else {
				buf = append(buf, '.')
			}
		}
		// east wall
		buf = append(buf, '#')
		ret = append(ret, buf)
	}
	// south wall
	ret = append(ret, bytes.Repeat([]byte{'#'}, 2*w+1))
	return ret
}

func (p *puzzle) distances(init pos) map[pos]int {
	type state struct {
		pos  pos
		dist int
	}
	work := []state{{
		pos:  init,
		dist: 0,
	}}
	seen := map[pos]int{}

	for i := 0; i < len(work); i++ {
		w := work[i]
		if v, ok := seen[w.pos]; ok {
			if v <= w.dist {
				continue
			}
		}
		seen[w.pos] = w.dist

		for d := N; d <= W; d++ {
			r := p.wat[w.pos]
			if r.Has(d) {
				work = append(work, state{
					pos:  w.pos.Walk(d),
					dist: w.dist + 1,
				})
			}
		}
	}

	return seen
}

func farthest(distances map[pos]int) int {
	ret := 0
	for _, v := range distances {
		ret = max(ret, v)
	}
	return ret
}

func atLeast1000(distances map[pos]int) int {
	ret := 0
	for _, v := range distances {
		if v >= 1000 {
			ret++
		}
	}
	return ret
}
