package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sample = `
199
200
208
210
200
207
240
269
260
263
`

func main1() {
	input := parse(sample)
	count := 0

	//for i := range input {
	for i := 0; i < len(input)-1; i++ {
		j := i + 1
		fmt.Print(i, input[i], input[j])
		if input[i] < input[j] {
			fmt.Println(" increase")
			count = count + 1
		} else if input[i] > input[j] {
			fmt.Println(" decrease")
		} else {
			fmt.Println(" same")
		}
	}

	fmt.Println(count)
}


func main() {
	input := parse(sample)
	count := 0

	last := 0
	for i := 2; i < len(input); i++ {
		sum := input[i-2] + input[i-1] + input[i]
		if i > 2 {
			fmt.Print(i, last, sum)
			if last < sum {
				fmt.Println(" increase")
				count = count + 1
			} else if last > sum {
				fmt.Println(" decrease")
			} else {
				fmt.Println(" same")
			}
		}
		last = sum

	}

	fmt.Println(count)
}

func parse(in string) []int {
	var ret []int
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			panic(fmt.Sprint(line, err))
		}
		ret = append(ret, n)
	}
	return ret
}
