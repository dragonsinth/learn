package main

import "fmt"

type fixedDie struct {
	rolls int
}

func (d *fixedDie) roll() int {
	ret := (d.rolls % 100) + 1
	d.rolls++
	return ret
}

func main1() {
	di := fixedDie{}

	p1pos, p2pos := 3, 7
	p1score, p2score := 0, 0

	for {
		p1pos += di.roll() + di.roll() + di.roll()
		p1pos %= 10
		p1score += p1pos + 1
		if p1score >= 1000 {
			break
		}

		p2pos += di.roll() + di.roll() + di.roll()
		p2pos %= 10
		p2score += p2pos + 1
		if p2score >= 1000 {
			break
		}
	}

	fmt.Println(p1score, p2score, di.rolls, p1score*di.rolls, p2score*di.rolls)
}
