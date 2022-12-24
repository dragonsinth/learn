package main

type allStates struct {
	max    pos
	states []state
}

func (s *allStates) get(turn int) *state {
	for len(s.states) <= turn {
		s.states = append(s.states, s.states[len(s.states)-1].next(s.max))
	}
	return &s.states[turn]
}

type state struct {
	blizzards   []blizzard
	blizzardPos map[pos][]blizzard
}

func (s state) next(max pos) state {
	bliz := make([]blizzard, len(s.blizzards))
	for i, b := range s.blizzards {
		pt := b.pt.advance(b.d)
		switch b.d {
		case NORTH:
			if pt.y == 0 {
				pt.y = max.y - 1
			}
		case SOUTH:
			if pt.y == max.y {
				pt.y = 1
			}
		case WEST:
			if pt.x == 0 {
				pt.x = max.x - 1
			}
		case EAST:
			if pt.x == max.x {
				pt.x = 1
			}
		default:
			panic(b.d)
		}
		bliz[i] = blizzard{pt: pt, d: b.d}
	}
	return state{
		blizzards:   bliz,
		blizzardPos: blizzardMap(bliz),
	}
}

func blizzardMap(bs []blizzard) map[pos][]blizzard {
	ret := make(map[pos][]blizzard, len(bs))
	for _, b := range bs {
		ret[b.pt] = append(ret[b.pt], b)
	}
	return ret
}

type blizzard struct {
	pt pos
	d  dir
}
