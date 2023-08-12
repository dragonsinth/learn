package grid

func Alloc2d[T any](w, h int, v T) [][]T {
	ret := make([][]T, h)
	buf := make([]T, w*h)

	for i := range buf {
		buf[i] = v
	}

	for y := 0; y < h; y++ {
		x := y * w
		ret[y] = buf[x : x+w]
	}
	return ret
}
