package main

import "github.com/dragonsinth/learn/aoc/2018/asm"

const data = `
#ip 4
seti 123 0 5
bani 5 456 5
eqri 5 72 5
addr 5 4 4
seti 0 0 4
seti 0 6 5
bori 5 65536 1
seti 4591209 6 5
bani 1 255 3
addr 5 3 5
bani 5 16777215 5
muli 5 65899 5
bani 5 16777215 5
gtir 256 1 3
addr 3 4 4
addi 4 1 4
seti 27 7 4
seti 0 0 3
addi 3 1 2
muli 2 256 2
gtrr 2 1 2
addr 2 4 4
addi 4 1 4
seti 25 4 4
addi 3 1 3
seti 17 0 4
setr 3 4 1
seti 7 2 4
eqrr 5 0 3
addr 3 4 4
seti 5 1 4
`

/*
 0 : fx = 123            // =123
 1 : fx = fx & 456       // =72
 2 : fx if fx == 72      // =1
 3 : pc = fx + pc        // GOTO 5
 4 : pc = 0              // GOTO 1
 5 : fx = 0              //
 6 : bx = fx | 65536     //
 7 : fx = 4591209        //
 8 : dx = bx & 255       //
 9 : fx = fx + dx        //
10 : fx = fx & 16777215  //
11 : fx = fx * 65899     //
12 : fx = fx & 16777215  //
13 : dx if 256 > bx      // [0 65536 0 0 13 12545763]
14 : pc = dx + pc        // JGT 16 -> 28
15 : pc = pc + 1         // JLE 17
16 : pc = 27             // GOTO 28
17 : dx = 0              //
18 : cx = dx + 1         // =1
19 : cx = cx * 256       // =256
20 : cx if cx > bx       //
21 : pc = cx + pc        // JGT 23 -> 26
22 : pc = pc + 1         // JLE 24
23 : pc = 25             // GOTO 26
24 : dx = dx + 1         //
25 : pc = 17             // GOTO 18
26 : bx = dx             //
27 : pc = 7              //
28 : dx if fx == ax      //
29 : pc = dx + pc        // JEQ 31 (HALT)
30 : pc = 5              // GOTO 6
*/

func translate(init int) asm.Registers {
	var _, bx, _, _, fx int

	for {
		bx = fx | 65536 // 0x10000
		fx = 4591209

		for {
			fx = fx + (bx & 255)
			fx = fx & 16777215 // 0xFFFFFF, 24 bits
			fx = fx * 65899
			fx = fx & 16777215 // 0xFFFFFF, 24 bits

			if bx < 256 {
				break
			}

			bx = bx / 256
		}

		if fx == init {
			return asm.Registers{init, bx, 1, 1, 30, fx}
		}
	}
}

func translate2() int {
	seen := map[int]int{}
	last := 0

	var _, bx, _, _, fx int

	for {
		bx = fx | 65536 // 0x10000
		fx = 4591209

		for {
			fx = fx + (bx & 255)
			fx = fx & 16777215 // 0xFFFFFF, 24 bits
			fx = fx * 65899
			fx = fx & 16777215 // 0xFFFFFF, 24 bits

			if bx < 256 {
				break
			}

			bx = bx / 256
		}

		if seen[fx] == 3 {
			return last
		}
		seen[fx]++
		last = fx
	}
}
