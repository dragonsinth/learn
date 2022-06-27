package main

import (
	"fmt"
	"strconv"
	"strings"
)

const sample = `
Player 1:
9
2
6
3
1

Player 2:
5
8
4
7
10
`

var (
	input = sample

	gameId = 1

	cachedResults = map[string]int{}
)

type puzzle struct {
	p1 []int
	p2 []int

	game  int
	round int

	log map[string]bool
}

func (p *puzzle) Key() string {
	return fmt.Sprintf("%+v/%+v", p.p1, p.p2)
}

func (p *puzzle) Run() int {
	fmt.Printf("=== Game %d ===\n", p.game)
	gameKey := p.Key()
	if w, ok := cachedResults[gameKey]; ok {
		fmt.Println()
		fmt.Printf("Player %d wins by cached result\n", w)
		return w
	}
	w := p.doRun()
	fmt.Println()
	fmt.Printf("The winner of game %d is player %d!\n", p.game, w)
	cachedResults[gameKey] = w
	return w
}

func (p *puzzle) doRun() int {
	for {
		stepKey := fmt.Sprintf("%+v/%+v", p.p1, p.p2)
		if p.log[stepKey] {
			fmt.Println()
			fmt.Println("Player 1 wins by duplication of state")
			return 1
		}
		p.log[stepKey] = true

		p.round++
		p.Step()
		if len(p.p1) == 0 {
			return 2
		}
		if len(p.p2) == 0 {
			return 1
		}
	}
}

func (p *puzzle) Step() {
	fmtPrintln()
	fmtPrintf("-- Round %d (Game %d) --\n", p.round, p.game)
	fmtPrintf("Player 1's deck: %+v\n", p.p1)
	fmtPrintf("Player 2's deck: %+v\n", p.p2)

	var c1, c2 int
	c1, p.p1 = p.p1[0], p.p1[1:]
	c2, p.p2 = p.p2[0], p.p2[1:]
	fmtPrintf("Player 1 plays %d\n", c1)
	fmtPrintf("Player 2 plays %d\n", c2)

	// See if we should play a subgame.
	winner := 0
	if c1 <= len(p.p1) && c2 <= len(p.p2) {
		fmt.Println("Playing a sub-game to determine the winner...")
		fmt.Println()
		subGame := puzzle{
			p1:    append([]int{}, p.p1[:c1]...),
			p2:    append([]int{}, p.p2[:c2]...),
			game:  gameId,
			round: 0,
			log:   map[string]bool{},
		}
		gameId++
		winner = subGame.Run()
		fmt.Printf("...anyway, back to game %d.\n", p.game)
	} else {
		if c1 < c2 {
			winner = 2
		} else if c2 < c1 {
			winner = 1
		} else {
			panic("same")
		}
	}

	fmtPrintf("Player %d wins round %d of game %d!\n", winner, p.round, p.game)

	switch winner {
	case 1:
		p.p1 = append(p.p1, c1, c2)
	case 2:
		p.p2 = append(p.p2, c2, c1)
	default:
		panic(winner)
	}
}

func (p puzzle) String() string {
	return ""
}

func main() {
	p := puzzle{
		p1:    nil,
		p2:    nil,
		game:  gameId,
		round: 0,
		log:   map[string]bool{},
	}
	gameId++

	state := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch line {
		case "Player 1:":
			state = 1
		case "Player 2:":
			state = 2
		default:
			v := mustInt(line)
			switch state {
			case 1:
				p.p1 = append(p.p1, v)
			case 2:
				p.p2 = append(p.p2, v)
			default:
				panic("no")
			}
		}
	}

	winner := p.Run()
	fmt.Println(p.round, winner)

	// Scoring
	cards := p.p1
	if len(cards) == 0 {
		cards = p.p2
	}
	fmt.Println(cards)
	score, mult := 0, len(cards)
	for _, v := range cards {
		score += v * mult
		mult--
	}
	fmt.Println(score)
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}

func fmtPrintln(args ...interface{}) {

}

func fmtPrintf(s string, args ...interface{}) {
}
