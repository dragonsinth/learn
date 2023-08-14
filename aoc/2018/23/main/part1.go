package main

type bot struct {
	p pos
	r int
}

func (b bot) toRegion() region {
	var ret region
	for i, v := range b.p {
		ret.min[i] = v - b.r
		ret.max[i] = v + b.r
	}
	return ret
}

func findPointsInBest(bots []bot) int {
	best, bestR := 0, 0
	for i, b := range bots {
		if b.r > bestR {
			best, bestR = i, b.r
		}
	}

	bestBot := bots[best]
	sum := 0
	for _, b := range bots {
		if dist(bestBot.p, b.p) <= bestBot.r {
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
