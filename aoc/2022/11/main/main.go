package main

import (
	"fmt"
	"sort"
)

type Monkey struct {
	items   []int
	op      func(old int) int
	test    int
	ifTrue  int
	ifFalse int
}

var sample = []Monkey{
	{
		items:   []int{79, 98},
		op:      func(old int) int { return old * 19 },
		test:    23,
		ifTrue:  2,
		ifFalse: 3,
	},
	{
		items:   []int{54, 65, 75, 74},
		op:      func(old int) int { return old + 6 },
		test:    19,
		ifTrue:  2,
		ifFalse: 0,
	},
	{
		items:   []int{79, 60, 97},
		op:      func(old int) int { return old * old },
		test:    13,
		ifTrue:  1,
		ifFalse: 3,
	},
	{
		items:   []int{74},
		op:      func(old int) int { return old + 3 },
		test:    17,
		ifTrue:  0,
		ifFalse: 1,
	},
}

func main() {
	run(clone(sample), false, 20)
	run(clone(sample), true, 10000)
}

func clone(monkeys []Monkey) []Monkey {
	var ret []Monkey
	for _, m := range monkeys {
		cp := m
		cp.items = append([]int{}, m.items...)
		ret = append(ret, cp)
	}
	return ret
}

func run(monkeys []Monkey, part2 bool, rounds int) {
	modulus := 1
	for _, monkey := range monkeys {
		modulus = modulus * monkey.test
	}
	fmt.Println(modulus)

	counters := make([]int, len(monkeys))
	for round := 1; round <= rounds; round++ {
		for i := range monkeys {
			monkey := &monkeys[i]
			fmtPrintf("Monkey %d:\n", i)

			for len(monkey.items) > 0 {
				// inspect
				item := monkey.items[0]
				fmtPrintf("  Monkey inspects an item with a worry level of %d.\n", item)
				counters[i]++
				monkey.items = monkey.items[1:]
				item = monkey.op(item)
				fmtPrintf("    Worry level is changed to %d.\n", item)

				// relief
				if part2 {
					item = item % modulus
				} else {
					item = item / 3
					fmtPrintf("    Monkey gets bored with item. Worry level is divided by 3 to %d.\n", item)
				}

				var which int
				if item%monkey.test == 0 {
					which = monkey.ifTrue
					fmtPrintf("    Current worry level is divisible by %d.\n", monkey.test)
				} else {
					which = monkey.ifFalse
					fmtPrintf("    Current worry level not is divisible by %d.\n", monkey.test)
				}
				fmtPrintf("    Item with worry level 500 is thrown to monkey %d.\n", which)
				target := &monkeys[which]
				target.items = append(target.items, item)
			}
		}

		if round == 1 || round == 20 || (round%1000 == 0) {
			fmt.Printf("\nRound: %d\n", round)
			for i, monkey := range monkeys {
				fmt.Printf("Monkey %d: %+v\n", i, monkey.items)
			}

			fmt.Printf("\n== After round %d ==\n", round)
			for i, counter := range counters {
				fmt.Printf("Monkey %d inspected items %d times\n", i, counter)
			}
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(counters)))
	fmt.Println(counters)
	fmt.Println(counters[0] * counters[1])
}

var debug = false

func fmtPrintf(format string, a ...any) {
	if debug {
		fmt.Printf(format, a...)
	}
}
