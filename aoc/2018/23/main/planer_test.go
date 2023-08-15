package main

import "testing"

func TestVerticesDiamond(t *testing.T) {
	bots := []bot{
		{p: pos{0, 0, 0}, rad: 1},
		{p: pos{5, 0, 0}, rad: 1},
		{p: pos{0, 5, 0}, rad: 1},
	}
	diamonds := toDiamonds(0, 1, bots)
	for _, d := range diamonds {
		t.Log(d.vertices())
	}
}
