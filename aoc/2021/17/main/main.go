package main

import "fmt"

var (
	sample = target{
		xMin: 20,
		xMax: 30,
		yMin: -10,
		yMax: -5,
	}

	data = target{
		xMin: 153,
		xMax: 199,
		yMin: -114,
		yMax: -75,
	}
)

type target struct {
	xMin, xMax, yMin, yMax int
}

type result int

func (r result) String() interface{} {
	switch r {
	case LOW:
		return "LOW"
	case HIGH:
		return "HIGH"
	case HIT:
		return "HIT"
	case PASS_THROUGH:
		return "PASS_THROUGH"
	default:
		panic(r)
	}
}

const (
	LOW = result(iota)
	HIGH
	HIT
	PASS_THROUGH
)

func main() {
	t := data

	// Start at x=1, y=0
	// calculate each launch as low, high, hit, pass through
	// increment y until it's too high, record max hit for this x
	// if there were any hits, we're done
	// increase x by 1 and try again

	for dx := 1; dx <= t.xMax; dx++ {
		if !isPossible(t, dx) {
			continue
		}

		bestY := computeBestY(t, dx)
		if bestY > 0 {
			fmt.Println(bestY)
			break
		}
	}

	hits := 0
	for dx := 1; dx <= t.xMax; dx++ {
		if !isPossible(t, dx) {
			continue
		}
		hits += computeHits(t, dx)
	}
	fmt.Println(hits)
}

func isPossible(t target, dx int) bool {
	// Is this shot even possible?
	x := 0
	for dx > 0 {
		x += dx
		dx--

		if x >= t.xMin && x <= t.xMax {
			return true
		}
	}
	return false
}

func computeBestY(t target, dx int) int {
	bestY := 0
	for dy := 1; dy <= -t.yMin; dy++ {
		maxY, result := computeTrajectory(t, dx, dy)
		fmt.Println(dx, dy, maxY, result.String())
		switch result {
		case LOW:
			continue
		case HIGH:
			return bestY
		case HIT:
			bestY = maxY
			continue
		case PASS_THROUGH:
			continue
		default:
			panic(result)
		}
	}
	return bestY
}

func computeHits(t target, dx int) int {
	hits := 0
	for dy := t.yMin; dy <= -t.yMin; dy++ {
		_, result := computeTrajectory(t, dx, dy)
		switch result {
		case LOW:
			continue
		case HIGH:
			return hits
		case HIT:
			hits++
			continue
		case PASS_THROUGH:
			continue
		default:
			panic(result)
		}
	}
	return hits
}

func computeTrajectory(t target, dx, dy int) (int, result) {
	maxY := 0
	x, y := 0, 0
	for {
		// move
		x += dx
		y += dy
		if dx > 0 {
			dx--
		}
		if y > maxY {
			maxY = y
		}
		dy--

		// check vs target

		if x < t.xMin && y < t.yMin {
			return 0, LOW // fell before reaching target
		}
		if x > t.xMax && y > t.yMax {
			return 0, HIGH // overshot
		}
		if x >= t.xMin && x <= t.xMax && y >= t.yMin && y <= t.yMax {
			return maxY, HIT // in the area
		}
		if x > t.xMax || y < t.yMin {
			return 0, PASS_THROUGH
		}
	}
}
