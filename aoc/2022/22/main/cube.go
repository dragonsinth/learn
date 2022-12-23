package main

type cube struct {
	data    [][]byte
	rsize   int
	regions [6]region
}

type region struct {
	id        int
	x, y      int
	translate [4]translation
}

type translation struct {
	dest int
	dir  dir
}

type cubePos struct {
	pos      pos
	regionId int
}

func (cp cubePos) advance(d dir) cubePos {
	return cubePos{
		pos:      cp.pos.advance(d),
		regionId: cp.regionId,
	}
}

type cubePlayer struct {
	cp cubePos
	d  dir
}

func (c *cube) to2d(cp cubePos) pos {
	r := c.regions[cp.regionId]
	p := cp.pos
	if p.x < 0 || p.x >= c.rsize {
		panic(p.x)
	}
	if p.y < 0 || p.y >= c.rsize {
		panic(p.y)
	}

	return pos{r.x + p.x, r.y + p.y}
}

func (c *cube) checkRegions() [][]byte {
	// verify translations
	saw := [6][4]bool{}
	for _, r := range c.regions {
		for _, t := range r.translate {
			if saw[t.dest][t.dir] {
				panic("here")
			}
			saw[t.dest][t.dir] = true
		}
	}
	for _, r := range saw {
		for _, t := range r {
			if !t {
				panic("here")
			}
		}
	}

	var out [][]byte
	// Compute bounds per row
	for y := 0; y < len(c.data); y++ {
		var buf []byte
		for x := 0; x < len(c.data[0]); x++ {
			b := c.data[y][x]
			if b == ' ' {
				buf = append(buf, ' ')
			} else if b != ' ' {
				// this better map to a region
				cpos := c.to3d(pos{x, y})
				buf = append(buf, byte(cpos.regionId+'0'))
			}
		}
		out = append(out, buf)
	}
	return out
}

func (c *cube) to3d(p pos) cubePos {
	// this better map to a region
	for _, r := range c.regions {
		x := p.x - r.x
		y := p.y - r.y
		if x < 0 || x >= c.rsize {
			continue
		}
		if y < 0 || y >= c.rsize {
			continue
		}
		return cubePos{
			pos:      pos{x, y},
			regionId: r.id,
		}
	}
	panic("no region")
}

func (c *cube) walk(player cubePlayer, dist int, onStep func()) cubePlayer {
	for dist > 0 {
		cp := player.cp.advance(player.d)
		pt := cp.pos
		var wat byte
		if pt.x >= 0 && pt.x < c.rsize && pt.y >= 0 && pt.y < c.rsize {
			// in bounds
			wat = c.get(c.to2d(cp))

			if wat == '#' {
				// wall, cannot move
				return player
			}

			// advance
			dist--
			player.cp = cp
			c.mark(player)
			onStep()
		} else {
			// out of bounds
			rg := c.regions[cp.regionId]
			tr := rg.translate[player.d]
			startDir := player.d
			endDir := tr.dir
			sz := c.rsize - 1
			for d := startDir; d != endDir; d = (d + 1) % 4 {
				// Rotate e.g x,y clockwise around center to align
				pt.x, pt.y = sz-pt.y, pt.x
			}

			// now warp to the opposite side
			switch endDir {
			case EAST:
				pt.x -= c.rsize
			case SOUTH:
				pt.y -= c.rsize
			case WEST:
				pt.x += c.rsize
			case NORTH:
				pt.y += c.rsize
			default:
				panic(endDir)
			}

			cp.pos = pt
			cp.regionId = tr.dest
			wat = c.get(c.to2d(cp))

			if wat == '#' {
				// wall, cannot move
				return player
			}

			// advance
			dist--
			player.cp = cp
			player.d = endDir
			c.mark(player)
			onStep()
		}
	}

	return player
}

func (c *cube) get(pt pos) byte {
	x, y := pt.x, pt.y
	if x < 0 || x >= len(c.data[0]) || y < 0 || y >= len(c.data) {
		return ' '
	}
	return c.data[y][x]
}

func (c *cube) mark(player cubePlayer) {
	pt := c.to2d(player.cp)
	c.data[pt.y][pt.x] = player.d.byte()
}
