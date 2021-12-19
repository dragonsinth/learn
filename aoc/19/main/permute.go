package main

func permutations(orig []point) [][]point {
	// Place each of 6 faces to be visible from -Z.
	faces := make([][]point, 6)
	faces[0] = orig                        // green is -Z, white is +Y
	faces[1] = xClockwise(faces[0])        // green -> white
	faces[2] = xClockwise(faces[1])        // white -> blue
	faces[3] = xClockwise(faces[2])        // blue -> yellow
	faces[4] = yClockwise(faces[0])        // green -> orange
	faces[5] = yCounterClockwise(faces[0]) // green -> red

	ret := make([][]point, 0, 24)
	for _, f := range faces {
		ret = append(ret, zRotations(f)...)
	}
	return ret
}

func zRotations(orig []point) [][]point {
	// Rotate the original cube 3 times clockwise around the Z axis
	ret := make([][]point, 4)
	ret[0] = orig
	for i := 1; i < 4; i++ {
		ret[i] = zClockwise(ret[i-1])
	}
	return ret
}

func xClockwise(points []point) []point {
	// Rotate clockwise while looking from -X
	// X := X, Y := Z, Z := -Y
	ret := make([]point, len(points))
	for i, p := range points {
		ret[i] = point{X: p[X], Y: p[Z], Z: -p[Y]}
	}
	return ret
}

func yClockwise(points []point) []point {
	// Rotate clockwise while looking from -Y
	// X := Z, Y := Y, Z := -X
	ret := make([]point, len(points))
	for i, p := range points {
		ret[i] = point{X: p[Z], Y: p[Y], Z: -p[X]}
	}
	return ret
}

func yCounterClockwise(points []point) []point {
	// Rotate counter-clockwise while looking from -Y
	// X := -Z, Y := Y, Z := X
	ret := make([]point, len(points))
	for i, p := range points {
		ret[i] = point{X: -p[Z], Y: p[Y], Z: p[X]}
	}
	return ret
}

func zClockwise(points []point) []point {
	// Rotate clockwise while looking from -Z
	// X := Y, Y := -X, Z := Z
	ret := make([]point, len(points))
	for i, p := range points {
		ret[i] = point{X: p[Y], Y: -p[X], Z: p[Z]}
	}
	return ret
}
