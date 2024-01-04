package main

type bot struct {
	p   pos
	rad int

	min pos
	max pos
}

func (b bot) inRange(p pos) bool {
	return dist(b.p, p) <= b.rad
}

func (b bot) normalize() bot {
	for d := range b.p {
		b.min[d] = b.p[d] - b.rad
		b.max[d] = b.p[d] + b.rad
	}
	return b
}

func findPointsInBest(bots []bot) int {
	best, bestR := 0, 0
	for i, b := range bots {
		if b.rad > bestR {
			best, bestR = i, b.rad
		}
	}

	bestBot := bots[best]
	sum := 0
	for _, b := range bots {
		if bestBot.inRange(b.p) {
			sum++
		}
	}
	return sum
}

func dist(p1 pos, p2 pos) int {
	var sum int
	for i := range p1 {
		sum += abs(p1[i] - p2[i])
	}
	return sum
}
