package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const sample = `
class: 1-3 or 5-7
row: 6-11 or 33-44
seat: 13-40 or 45-50

your ticket:
7,1,14

nearby tickets:
7,3,47
40,4,50
55,2,20
38,6,12
`

var (
	input = sample

	ruleLine = regexp.MustCompile(`^([^:]+): (\d+)-(\d+) or (\d+)-(\d+)`)
)

type rule struct {
	name                   string
	min1, max1, min2, max2 int
}

func (r rule) matches(v int) bool {
	return (r.min1 <= v && v <= r.max1) || (r.min2 <= v && v <= r.max2)
}

func main() {
	parseState := 0
	var myTicket []int
	var rules []rule
	var tickets [][]int
	errorRate := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch parseState {
		case 0:
			if line == "your ticket:" {
				parseState++
				continue
			}
			if !ruleLine.MatchString(line) {
				panic(line)
			}
			matches := ruleLine.FindStringSubmatch(line)
			if len(matches) != 6 {
				panic(matches)
			}

			r := rule{
				name: matches[1],
				min1: mustInt(matches[2]),
				max1: mustInt(matches[3]),
				min2: mustInt(matches[4]),
				max2: mustInt(matches[5]),
			}
			rules = append(rules, r)
		case 1:
			if line == "nearby tickets:" {
				parseState++
				continue
			}
			myTicket = parseTicket(line)
		case 2:
			ticket := parseTicket(line)
			if len(ticket) != len(myTicket) {
				panic(ticket)
			}
			valid := true
			for _, v := range ticket {
				anyMatch := false
				for _, r := range rules {
					if r.matches(v) {
						anyMatch = true
						break
					}
				}
				if !anyMatch {
					errorRate += v
					valid = false
				}
			}
			if valid {
				tickets = append(tickets, ticket)
			}
		}
	}

	if parseState != 2 {
		panic(parseState)
	}
	fmt.Println(myTicket)
	fmt.Println(errorRate)

	// Map the rules to possibilities.
	possible := map[string][]bool{}
	for _, r := range rules {
		possible[r.name] = make([]bool, len(myTicket))
		for pos := range myTicket {
			match := true
			for _, ticket := range tickets {
				if !r.matches(ticket[pos]) {
					match = false
					break
				}
			}
			if match {
				possible[r.name][pos] = true
			}
		}
	}

	rulePos := map[string]int{}
	for len(rulePos) < len(myTicket) {
		// Find a rule with only one possibility.
		found := false
		for _, r := range rules {
			if _, ok := rulePos[r.name]; ok {
				continue // already solved this rule
			}
			flags := possible[r.name]
			count := 0
			pos := 0
			for p, f := range flags {
				if f {
					count++
					pos = p
				}
			}
			if count == 0 {
				panic(r.name)
			}
			if count == 1 {
				//fmt.Println(r.name, pos)
				found = true
				rulePos[r.name] = pos
				// No one else can have this position
				for _, v := range possible {
					v[pos] = false
				}
				break
			}
			// keep going
		}
		if !found {
			panic("none found")
		}
	}

	// Now map myTicket
	prod := 1
	for k, v := range rulePos {
		fmt.Println(k, myTicket[v])
		if strings.HasPrefix(k, "departure") {
			prod = prod * myTicket[v]
		}
	}
	fmt.Println(prod)
}

func parseTicket(line string) []int {
	var ret []int
	for _, part := range strings.Split(line, ",") {
		ret = append(ret, mustInt(part))
	}
	return ret
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
