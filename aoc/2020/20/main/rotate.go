package main

func rotate(in [][]bool, or orient) [][]bool {
	ysz, xsz := len(in), len(in[0])
	xInc, yInc := 0, 0
	flip := false
	switch or {
	case NORT:
		xInc, yInc = 1, 1
		flip = false
	case EAST:
		xInc, yInc = 1, -1
		flip = true
	case SOUT:
		xInc, yInc = -1, -1
		flip = false
	case WEST:
		xInc, yInc = -1, 1
		flip = true

	case TRON:
		xInc, yInc = -1, 1
		flip = false
	case TSAE:
		xInc, yInc = -1, -1
		flip = true
	case TUOS:
		xInc, yInc = 1, -1
		flip = false
	case TSEW:
		xInc, yInc = 1, 1
		flip = true
	default:
		panic(or)
	}

	ymin, ymax := 0, ysz
	if yInc < 0 {
		ymin, ymax = ysz-1, -1
	}
	xmin, xmax := 0, xsz
	if xInc < 0 {
		xmin, xmax = xsz-1, -1
	}

	if !flip {
		ret := makeBoolField(ysz, xsz)
		wy := 0
		for y := ymin; y != ymax; y += yInc {
			wx := 0
			for x := xmin; x != xmax; x += xInc {
				ret[wy][wx] = in[y][x]
				wx++
			}
			wy++
		}
		return ret
	} else {
		ret := makeBoolField(xsz, ysz)
		wy := 0
		for x := xmin; x != xmax; x += xInc {
			wx := 0
			for y := ymin; y != ymax; y += yInc {
				ret[wy][wx] = in[y][x]
				wx++
			}
			wy++
		}
		return ret
	}
}

func toString(in [][]bool) string {
	var buf []byte
	for _, line := range in {
		for _, v := range line {
			if v {
				buf = append(buf, '#')
			} else {
				buf = append(buf, '.')
			}
		}
		buf = append(buf, '\n')
	}
	return string(buf)
}

func invert(want uint16) uint16 {
	ret := uint16(0)
	for i := 0; i < size; i++ {
		ret <<= 1
		if want&1 == 1 {
			ret |= 1
		}
		want >>= 1
	}
	return ret
}
