package asm

import (
	"fmt"
	"strconv"
	"strings"
)

func Parse(input string) *Prog {
	p := Prog{
		IpReg: -1,
		Names: [6]string{
			"ax",
			"bx",
			"cx",
			"dx",
			"ex",
			"fx",
		},
	}

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if p.IpReg < 0 {
			if !dirRe.MatchString(line) {
				panic(line)
			}
			m := dirRe.FindStringSubmatch(line)
			p.IpReg = mustInt(m[1])
			p.Names[p.IpReg] = "pc"
		} else {
			if !instRe.MatchString(line) {
				panic(line)
			}
			m := instRe.FindStringSubmatch(line)
			name := m[1]
			f := Instrs[name]
			if f == nil {
				panic(name)
			}
			p.Code = append(p.Code, Instr{
				Name: name,
				Fun:  f,
				Ops: [3]int{
					mustInt(m[2]),
					mustInt(m[3]),
					mustInt(m[4]),
				},
			})

		}
	}
	return &p
}

func mustInt(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		panic(fmt.Sprint(s, err))
	} else {
		return v
	}
}
