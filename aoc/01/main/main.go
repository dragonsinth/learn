package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
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

func parse(in string) []int {
	var ret []int
	split := strings.Split(in, "\n")
	for _, s := range split {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(fmt.Sprint(s, err))
		}
		ret = append(ret, n)
	}
	return ret
}
