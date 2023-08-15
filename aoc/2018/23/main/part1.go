package main

type bot struct {
	p   pos
	rad int
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
		if dist(bestBot.p, b.p) <= bestBot.rad {
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
