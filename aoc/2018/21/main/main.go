package main

import (
	"fmt"
	"github.com/dragonsinth/learn/aoc/2018/asm"
)

func main() {
	p := asm.Parse(data)
	p.Decompile()
	p.Run(true, 12446070)
	fmt.Println(translate(12446070))
	fmt.Println(translate2())
}
