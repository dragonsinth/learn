package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"testing"
)

const (
	hdr = `package main

func trial(x, y [63]bool) [63]bool {
	var z [63]bool
`

	ftr = `
	return z
}
`
)

// solution here (redacted)
var swaps = map[string]string{}

func TestAnswer(t *testing.T) {
	var ret []string
	for k := range swaps {
		ret = append(ret, k)
	}
	slices.Sort(ret)
	t.Logf(strings.Join(ret, ","))
}

// REDACTED
const dataGates = `
x00 XOR y00 -> z00

x01 XOR y01 -> aaa
x00 AND y00 -> bbb
aaa XOR bbb -> z01

x02 XOR y02 -> ccc
x01 AND y01 -> ddd
aaa AND bbb -> eee
eee OR ddd -> fff
ccc XOR fff -> z02
`

func TestGen(t *testing.T) {
	gates := parseGates(dataGates)

	for i := range gates {
		g := &gates[i]
		if swp := swaps[g.o]; swp != "" {
			g.o = swp
		}
	}

	slices.SortFunc(gates, func(a, b gate) int {
		return strings.Compare(a.o, b.o) // sort up
	})

	gatesByOutput := map[string]gate{}
	for _, g := range gates {
		gatesByOutput[g.o] = g
	}

	var w bytes.Buffer
	w.WriteString(hdr)

	seen := map[string]bool{}
	for _, g := range gates {
		output := g.o
		if output[0] == 'z' {
			_, _ = fmt.Fprintf(&w, "\n\t// %s\n", output)
			printTransitiveClosure(&w, seen, gatesByOutput, output)
		}
	}

	w.WriteString(ftr)

	_ = os.WriteFile("main.gen.go", w.Bytes(), 0644)
}

var sym = []string{"&&", "||", "!="}

func printTransitiveClosure(w io.Writer, seen map[string]bool, output map[string]gate, name string) {
	if seen[name] {
		return
	}
	seen[name] = true

	if g, ok := output[name]; ok {
		printTransitiveClosure(w, seen, output, g.a)
		printTransitiveClosure(w, seen, output, g.b)
		_, _ = fmt.Fprintf(w, "\t%s %s %s %s %s\n", regName(name), eq(name), regName(g.a), sym[g.op], regName(g.b))
	}
}

func eq(name string) string {
	if name[0] >= 'x' && name[0] <= 'z' {
		return "="
	} else {
		return ":="
	}
}

func regName(name string) string {
	if name[0] >= 'x' && name[0] <= 'z' {
		num := int(name[1]-'0')*10 + int(name[2]-'0')
		return fmt.Sprintf("%c[%d]", name[0], num)
	} else {
		return name
	}
}
