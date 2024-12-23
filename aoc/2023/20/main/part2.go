package main

import (
	"fmt"
	"sort"
	"strings"
)

func (p puz) memDump(mem memory, step int) (string, bool) {
	ok := true
	var sb strings.Builder
	keys := mapKeys(p.nodes)
	sort.Strings(keys)
	for _, k := range keys {
		p := mem.vals[k]
		if exfn := expect[k]; exfn != nil {
			expv := exfn(step)
			if p != expv {
				fmt.Println("X", k, step, expv, p)
				ok = false
			}
		}
		v := 0
		if p {
			v = 1
		}
		fmt.Fprintf(&sb, "%s:%d,", k, v)
	}
	return sb.String(), ok
}

func modShiftSub(shift int, sub int) func(int) pulse {
	return func(i int) pulse {
		i = i % sub
		return (i >> shift % 2) == 1
	}
}

var expect = map[string]func(int) pulse{
	start: func(i int) pulse {
		return LO
	},

	"fv": modShiftSub(0, 3739),
	"ls": modShiftSub(0, 3821),
	"ts": modShiftSub(0, 3889),
	"bn": modShiftSub(0, 4093),
	"cd": modShiftSub(1, 3739),
	"tb": modShiftSub(1, 3821),
	"lp": modShiftSub(1, 3889),
	"xz": modShiftSub(1, 4093),
	"tc": modShiftSub(2, 3739),
	"ss": modShiftSub(2, 3821),
	"xx": modShiftSub(2, 3889),
	"dm": modShiftSub(2, 4093),
	"vp": modShiftSub(3, 3739),
	"mb": modShiftSub(3, 3821),
	"zd": modShiftSub(3, 3889),
	"ct": modShiftSub(3, 4093),
	"jb": modShiftSub(4, 3739),
	"hv": modShiftSub(4, 3821),
	"kg": modShiftSub(4, 3889),
	"rj": modShiftSub(4, 4093),
	"hc": modShiftSub(5, 3739),
	"nh": modShiftSub(5, 3821),
	"gd": modShiftSub(5, 3889),
	"qk": modShiftSub(5, 4093),
	"kh": modShiftSub(6, 3739),
	"xm": modShiftSub(6, 3821),
	"bs": modShiftSub(6, 3889),
	"xg": modShiftSub(6, 4093),
	"ck": modShiftSub(7, 3739),
	"dp": modShiftSub(7, 3821),
	"gk": modShiftSub(7, 3889),
	"pb": modShiftSub(7, 4093),
	"vg": modShiftSub(8, 3739),
	"ql": modShiftSub(8, 3821),
	"pc": modShiftSub(8, 3889),
	"rz": modShiftSub(8, 4093),
	"mm": modShiftSub(9, 3739),
	"tr": modShiftSub(9, 3821),
	"hr": modShiftSub(9, 3889),
	"cr": modShiftSub(9, 4093),
	"jr": modShiftSub(10, 3739),
	"hl": modShiftSub(10, 3821),
	"vf": modShiftSub(10, 3889),
	"rf": modShiftSub(10, 4093),
	"rb": modShiftSub(11, 3739),
	"vj": modShiftSub(11, 3821),
	"nc": modShiftSub(11, 3889),
	"xk": modShiftSub(11, 4093),
}
