package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var sample = `
939
7,13,x,x,59,x,31,19
`

var (
	input = sample
)

func main1() {
	ln := 0
	var depart int
	var busses []int
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch ln {
		case 0:
			depart = mustInt(line)
		case 1:
			parts := strings.Split(line, ",")
			for _, p := range parts {
				if p == "x" {
					continue
				}
				busses = append(busses, mustInt(p))
			}
		default:
			panic(line)
		}

		ln++
	}

	type trip struct {
		bus  int
		when int
	}
	var trips []trip
	for _, bus := range busses {
		n := depart / bus
		when := bus * (n + 1)
		fmt.Println(bus, when)
		trips = append(trips, trip{
			bus:  bus,
			when: when,
		})
	}

	sort.Slice(trips, func(i, j int) bool {
		return trips[i].when < trips[j].when
	})

	bus, when := trips[0].bus, trips[0].when
	fmt.Println(bus, when, when-depart, bus*(when-depart))
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
