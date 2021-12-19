package main

import (
	"fmt"
	"strconv"
	"strings"
)

func parseScanners(input string) []*scanner {
	id := -1
	var scanners []*scanner
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if scannerLine.MatchString(line) {
			subs := scannerLine.FindStringSubmatch(line)
			which := mustInt(subs[1])
			id++
			if which != id {
				panic(which)
			}
			scanners = append(scanners, &scanner{id: id})
			continue
		}

		parts := strings.Split(line, ",")
		x, y, z := mustInt(parts[0]), mustInt(parts[1]), mustInt(parts[2])
		scanners[id].points = append(scanners[id].points, point{x, y, z})
	}
	return scanners
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
