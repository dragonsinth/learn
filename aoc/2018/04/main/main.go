package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const sample = `
[1518-11-01 00:00] Guard #10 begins shift
[1518-11-01 00:05] falls asleep
[1518-11-01 00:25] wakes up
[1518-11-01 00:30] falls asleep
[1518-11-01 00:55] wakes up
[1518-11-01 23:58] Guard #99 begins shift
[1518-11-02 00:40] falls asleep
[1518-11-02 00:50] wakes up
[1518-11-03 00:05] Guard #10 begins shift
[1518-11-03 00:24] falls asleep
[1518-11-03 00:29] wakes up
[1518-11-04 00:02] Guard #99 begins shift
[1518-11-04 00:36] falls asleep
[1518-11-04 00:46] wakes up
[1518-11-05 00:03] Guard #99 begins shift
[1518-11-05 00:45] falls asleep
[1518-11-05 00:55] wakes up
`

var re = regexp.MustCompile(`^Guard #(\d+) begins shift$`)

func main() {
	track := map[int][]int{}
	var curId int
	var curDayAsleep []bool
	var sleepyTime int

	finishShift := func() {
		if sleepyTime >= 0 {
			panic("guard asleep at end of shift")
		}
		if track[curId] == nil {
			track[curId] = make([]int, 60)
		}
		for i, v := range curDayAsleep {
			if v {
				track[curId][i]++
			}
		}
	}

	lines := strings.Split(sample, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "] ")
		tm, err := time.Parse(`2006-01-02 15:04`, parts[0][1:])
		if err != nil {
			panic(err)
		}
		now := tm.Minute()

		s := parts[1]
		switch s {
		case "falls asleep":
			if curId == 0 {
				panic("no guard")
			}
			if sleepyTime >= 0 {
				panic("guard was already asleep")
			}
			sleepyTime = now
		case "wakes up":
			if curId == 0 {
				panic("no guard")
			}
			if sleepyTime < 0 {
				panic("guard was not asleep")
			}
			// mark the sleep duration
			for i := sleepyTime; i < now; i++ {
				curDayAsleep[i] = true
			}
			sleepyTime = -1
		default:
			if !re.MatchString(s) {
				panic(s)
			}
			m := re.FindStringSubmatch(s)
			id := mustInt(m[1])

			// Finish the previous shift
			if curId != 0 {
				finishShift()
			}

			// prepare the next shift
			curId = id
			curDayAsleep = make([]bool, 60)
			sleepyTime = -1
		}
	}
	if curId != 0 {
		finishShift()
	}

	for k, v := range track {
		fmt.Println(k, v)
	}

	bestId := findBestId(track)
	bestMinute := findBestMinute(track[bestId])
	fmt.Println(bestId, bestMinute, bestId*bestMinute)

	// part 2
	bestId, bestMinute = findBestGuardMinute(track)
	fmt.Println(bestId, bestMinute, bestId*bestMinute)

}

func findBestId(track map[int][]int) int {
	var bestId, bestSum int
	for id, row := range track {
		sum := 0
		for _, c := range row {
			sum += c
		}
		if sum > bestSum {
			bestId, bestSum = id, sum
		}
	}
	return bestId
}

func findBestMinute(row []int) int {
	var best, bestV int
	for i, v := range row {
		if v > bestV {
			best, bestV = i, v
		}
	}
	return best
}

func findBestGuardMinute(track map[int][]int) (int, int) {
	var bestId, bestMin, bestScore int
	for id, row := range track {
		for i, c := range row {
			if c > bestScore {
				bestId, bestMin, bestScore = id, i, c
			}
		}
	}
	return bestId, bestMin
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
