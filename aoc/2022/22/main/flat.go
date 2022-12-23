package main

type flat struct {
	data          [][]byte
	width, height int

	xMins []int
	xMaxs []int
	yMins []int
	yMaxs []int
}

func (f *flat) computeBounds() {
	// Compute bounds per row
	for y := 0; y < f.height; y++ {
		xMin, xMax := f.width, 0
		for x := 0; x < f.width; x++ {
			if f.data[y][x] != ' ' {
				xMin = min(xMin, x)
				xMax = max(xMax, x)
			}
		}
		f.xMins = append(f.xMins, xMin)
		f.xMaxs = append(f.xMaxs, xMax)
	}
	// Compute bounds per col
	for x := 0; x < f.width; x++ {
		yMin, yMax := f.width, 0
		for y := 0; y < f.height; y++ {
			if f.data[y][x] != ' ' {
				yMin = min(yMin, y)
				yMax = max(yMax, y)
			}
		}
		f.yMins = append(f.yMins, yMin)
		f.yMaxs = append(f.yMaxs, yMax)
	}
}

func (f *flat) walk(player flatPlayer, dist int, onStep func()) flatPlayer {
	for dist > 0 {
		pos := player.pos.advance(player.d)

		wat := f.get(pos.x, pos.y)
		if wat == ' ' {
			// warp around
			switch player.d {
			case EAST:
				pos.x = f.xMins[pos.y]
			case SOUTH:
				pos.y = f.yMins[pos.x]
			case WEST:
				pos.x = f.xMaxs[pos.y]
			case NORTH:
				pos.y = f.yMaxs[pos.x]
			default:
				panic(player.d)
			}
			wat = f.get(pos.x, pos.y)
			if wat == ' ' {
				panic("should not be blank")
			}
		}

		if wat == '#' {
			// wall, cannot move
			return player
		}

		// advance
		dist--
		player.x, player.y = pos.x, pos.y
		f.mark(player)
		onStep()
	}

	return player
}

func (f *flat) get(x int, y int) byte {
	if x < 0 || x >= f.width || y < 0 || y >= f.height {
		return ' '
	}
	return f.data[y][x]
}

func (f *flat) mark(player flatPlayer) {
	f.data[player.y][player.x] = player.d.byte()
}
