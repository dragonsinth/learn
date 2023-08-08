package main

import (
	"github.com/dragonsinth/learn/aoc/2018/asm"
	"golang.org/x/exp/slices"
)

const data = `
#ip 5
addi 5 16 5
seti 1 0 4
seti 1 8 1
mulr 4 1 3
eqrr 3 2 3
addr 3 5 5
addi 5 1 5
addr 4 0 0
addi 1 1 1
gtrr 1 2 3
addr 5 3 5
seti 2 4 5
addi 4 1 4
gtrr 4 2 3
addr 3 5 5
seti 1 7 5
mulr 5 5 5
addi 2 2 2
mulr 2 2 2
mulr 5 2 2
muli 2 11 2
addi 3 6 3
mulr 3 5 3
addi 3 9 3
addr 2 3 2
addr 5 0 5
seti 0 5 5
setr 5 9 3
mulr 3 5 3
addr 5 3 3
mulr 5 3 3
muli 3 14 3
mulr 3 5 3
addr 2 3 2
seti 0 1 0
seti 0 0 5
`

/*
 0 : pc = pc + 16        ; GOTO 17
 1 : ex = 1              ;
 2 : bx = 1              ;
 3 : dx = ex * bx        ;
 4 : dx if dx == cx      ;
 5 : pc = dx + pc        ; JEQ 7
 6 : pc = pc + 1         ; JNE 8
 7 : ax = ex + ax        ;
 8 : bx = bx + 1         ;
 9 : dx if bx > cx       ;
10 : pc = pc + dx        ; JGT 12
11 : pc = 2              ; JLE 3
12 : ex = ex + 1         ;
13 : dx if ex > cx       ;
14 : pc = dx + pc        ; JGT 16 (HALT)
15 : pc = 1              ; JLE 2
16 : pc = pc * pc        ; HALT
17 : cx = cx + 2         ;
18 : cx = cx * cx        ;
19 : cx = 19 * cx        ;
20 : cx = cx * 11        ;
21 : dx = dx + 6         ;
22 : dx = dx * 22        ;
23 : dx = dx + 9         ;
24 : cx = cx + dx        ;
25 : pc = pc + ax        ; if ax == 1, GOTO 27; if ax big, HALT?
26 : pc = 0              ; if ax == 0, GOTO 1
27 : dx = 27             ;
28 : dx = dx * 28        ;
29 : dx = 29 + dx        ;
30 : dx = 30 * dx        ;
31 : dx = dx * 14        ;
32 : dx = dx * 32        ;
33 : cx = cx + dx        ;
34 : ax = 0              ;
35 : pc = 0              ; GOTO 1
*/

func translate(init int) asm.Registers {
	var ax, bx, cx, dx, ex int
	ax = init

	goto _17

_1:
	ex = 1
_2:
	bx = 1
_3:
	dx = ex * bx

	if dx == cx {
		ax = ex + ax
	}

	bx = bx + 1

	if bx <= cx {
		goto _3
	}

	ex = ex + 1

	if ex > cx {
		return asm.Registers{ax, bx, cx, dx, ex, 999}
	} else {
		goto _2
	}

_17:
	cx = cx + 2
	cx = cx * cx
	cx = 19 * cx
	cx = cx * 11
	dx = dx + 6
	dx = dx * 22
	dx = dx + 9
	cx = cx + dx

	switch ax {
	case 0:
	case 1:
		dx = 27
		dx = dx * 28
		dx = 29 + dx
		dx = 30 * dx
		dx = dx * 14
		dx = dx * 32
		cx = cx + dx
		ax = 0
	default:
		panic(ax)
	}

	goto _1
}

func primeFactors(n int) []int {
	var fac []int
	for i := 2; i*i <= n; i++ {
		for n%i == 0 {
			fac = append(fac, i)
			n = n / i
		}
	}
	if n > 1 {
		fac = append(fac, n)
	}
	return fac
}

func combine(fac []int) []int {
	var ret []int

	for i, c := 0, 1<<len(fac); i < c; i++ {
		prod := 1
		for j, v := range fac {
			if i&(1<<j) != 0 {
				prod *= v
			}
		}
		ret = append(ret, prod)
	}

	slices.Sort(ret)
	slices.Compact(ret)
	return ret
}

func sum(in []int) int {
	ret := 0
	for _, n := range in {
		ret += n
	}
	return ret
}
