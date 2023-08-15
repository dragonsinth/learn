package main

func zipperMerge(a []int, b []int) []int {
	r := make([]int, len(a)+len(b))
	ai, bi, ri := 0, 0, 0
	for ai < len(a) && bi < len(b) {
		av := a[ai]
		bv := b[bi]
		switch {
		case av < bv:
			r[ri] = av
			ri++
			ai++
		case bv < av:
			r[ri] = bv
			ri++
			bi++
		case av == bv:
			r[ri] = av
			ri++
			ai++
			bi++
		default:
			panic(av)
		}
	}
	for ai < len(a) {
		r[ri] = a[ai]
		ri++
		ai++
	}
	for bi < len(b) {
		r[ri] = b[bi]
		ri++
		bi++
	}
	return r[:ri]
}
