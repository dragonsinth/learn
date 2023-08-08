package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/2018/asm"
)

const sample = `
#ip 0
seti 5 0 1
seti 6 0 2
addi 0 1 0
addr 1 2 3
setr 1 0 0
seti 8 0 4
seti 9 0 5
`

func main() {
	p := asm.Parse(sample)
	p.Decompile()
	p.Run(true, 0)

	p = asm.Parse(data)
	p.Decompile()
	p.Run(true, 0)

	// 977, 10551377
	reg := translate(0)
	fmt.Println(reg)

	for _, n := range []int{977, 10551377} {
		pf := primeFactors(n)
		fmt.Println(pf)
		f := combine(pf)
		fmt.Println(f)
		fmt.Println(sum(f))
	}
}
