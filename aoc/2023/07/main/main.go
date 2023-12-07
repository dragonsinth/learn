package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

var sample = `
32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483
`

type hand struct {
	cards []byte
	rank  rank
	bid   int
}

func (h hand) String() string {
	return fmt.Sprintf("%s, %d = %d", string(h.cards), h.rank, h.bid)
}

func rankSort(a, b hand) bool {
	if a.rank != b.rank {
		return a.rank < b.rank
	}
	return bytes.Compare(a.cards, b.cards) < 0
}

type rank int

const (
	HIGH = rank(iota)
	ONE
	TWO
	THREE
	FULL
	FOUR
	FIVE
)

func mapCards(in string, m []byte) []byte {
	if len(in) != 5 {
		panic(in)
	}
	cards := make([]byte, 5)
	for i, c := range in {
		if m[c] == 0 {
			panic(c)
		}
		cards[i] = m[c]
	}
	return cards
}

func main() {
	part1(sample)
	part2(sample)
}

func part1(input string) {
	var hands []hand
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			panic(parts)
		}
		bid := mustInt(parts[1])

		cards := mapCards(parts[0], []byte{
			'2': 'a',
			'3': 'b',
			'4': 'c',
			'5': 'd',
			'6': 'e',
			'7': 'f',
			'8': 'g',
			'9': 'h',
			'T': 'i',
			'J': 'j',
			'Q': 'k',
			'K': 'l',
			'A': 'm',
		})
		hands = append(hands, hand{
			cards: cards,
			rank:  computeRank(cards),
			bid:   bid,
		})
	}

	slices.SortFunc(hands, rankSort)

	sum := 0
	for i, h := range hands {
		sum += (i + 1) * h.bid
		fmt.Println(h)
	}
	fmt.Println(sum)
}

func part2(input string) {
	var hands []hand
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			panic(parts)
		}
		bid := mustInt(parts[1])

		cards := mapCards(parts[0], []byte{
			'J': 'a',
			'2': 'b',
			'3': 'c',
			'4': 'd',
			'5': 'e',
			'6': 'f',
			'7': 'g',
			'8': 'h',
			'9': 'i',
			'T': 'j',
			'Q': 'k',
			'K': 'l',
			'A': 'm',
		})
		hands = append(hands, hand{
			cards: cards,
			rank:  computeRankWithJokers(cards),
			bid:   bid,
		})

	}

	slices.SortFunc(hands, rankSort)

	sum := 0
	for i, h := range hands {
		sum += (i + 1) * h.bid
		fmt.Println(h)
	}
	fmt.Println(sum)
}

func countCards(cards []byte) []int {
	counts := make([]int, 13)
	for _, c := range cards {
		counts[c-'a']++
	}
	return counts
}

func scoreCounts(counts []int) rank {
	hasPair := false
	hasTrio := false
	for _, count := range counts {
		switch count {
		case 2:
			if hasTrio {
				return FULL
			} else if hasPair {
				return TWO
			} else {
				hasPair = true
			}
		case 3:
			if hasTrio {
				panic("here")
			} else if hasPair {
				return FULL
			} else {
				hasTrio = true
			}
		case 4:
			return FOUR
		case 5:
			return FIVE
		}
	}
	if hasTrio {
		return THREE
	} else if hasPair {
		return ONE
	} else {
		return HIGH
	}
}

func computeRank(cards []byte) rank {
	counts := countCards(cards)
	return scoreCounts(counts)
}

func computeRankWithJokers(cards []byte) rank {
	counts := countCards(cards)

	// Transfer the jokers to whatever already has the highest.
	best, bestIdx := 0, 0
	numJokers := counts[0]
	for i, count := range counts {
		if i > 0 && count > best {
			bestIdx, best = i, count
		}
	}
	counts[bestIdx] += numJokers
	counts[0] -= numJokers

	return scoreCounts(counts)
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
