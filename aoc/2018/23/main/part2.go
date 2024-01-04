package main

import (
	"math/rand"
)

func calcScore(p pos, bots []bot) int {
	ret := 0
	for _, b := range bots {
		if b.inRange(p) {
			ret++
		}
	}
	return ret
}

func randomMatch(search pos, rad pos, iter int, bots []bot) (int, pos) {
	bestScore, bestPos := calcScore(search, bots), search
	for i := 0; i < iter; i++ {
		var p pos
		for d := X; d <= Z; d++ {
			r := rad[d]
			p[d] = search[d] + rand.Intn(r*2+1) - r
		}
		score := calcScore(p, bots)
		if score > bestScore || (score == bestScore && p.rad() < bestPos.rad()) {
			bestScore, bestPos = score, p
		}
	}
	return bestScore, bestPos
}

func exhaustiveMatch(search pos, rad int, bots []bot) (int, pos) {
	bestScore, bestPos := calcScore(search, bots), search
	for x := search[X] - rad; x <= search[X]+rad; x++ {
		for y := search[Y] - rad; y <= search[Y]+rad; y++ {
			for z := search[Z] - rad; z <= search[Z]+rad; z++ {
				p := pos{x, y, z}
				score := calcScore(p, bots)
				if score > bestScore || (score == bestScore && p.rad() < bestPos.rad()) {
					bestScore, bestPos = score, p
				}
			}
		}
	}
	return bestScore, bestPos
}
