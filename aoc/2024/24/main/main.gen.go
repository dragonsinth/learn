package main

func trial(x, y [63]bool) [63]bool {
	var z [63]bool

	// z00
	z[0] = x[0] != y[0]

	// z01
	aaa := x[1] != y[1]
	bbb := x[0] && y[0]
	z[1] = aaa != bbb

	// z02
	ccc := x[2] != y[2]
	eee := aaa && bbb
	ddd := x[1] && y[1]
	fff := eee || ddd
	z[2] = ccc != fff

	return z
}
