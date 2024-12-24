package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

const sample1Regs = `
x00: 1
x01: 1
x02: 1
y00: 0
y01: 1
y02: 0
`

const sample1Gates = `
x00 AND y00 -> z00
x01 XOR y01 -> z01
x02 OR y02 -> z02
`

const sample2Regs = `
x00: 1
x01: 0
x02: 1
x03: 1
x04: 0
y00: 1
y01: 1
y02: 1
y03: 1
y04: 1
`

const sample2Gates = `
ntg XOR fgs -> mjb
y02 OR x01 -> tnw
kwq OR kpj -> z05
x00 OR x03 -> fst
tgd XOR rvg -> z01
vdt OR tnw -> bfw
bfw AND frj -> z10
ffh OR nrd -> bqk
y00 AND y03 -> djm
y03 OR y00 -> psh
bqk OR frj -> z08
tnw OR fst -> frj
gnj AND tgd -> z11
bfw XOR mjb -> z00
x03 OR x00 -> vdt
gnj AND wpb -> z02
x04 AND y00 -> kjc
djm OR pbm -> qhw
nrd AND vdt -> hwm
kjc AND fst -> rvg
y04 OR y02 -> fgs
y01 AND x02 -> pbm
ntg OR kjc -> kwq
psh XOR fgs -> tgd
qhw XOR tgd -> z09
pbm OR djm -> kpj
x03 XOR y03 -> ffh
x00 XOR y04 -> ntg
bfw OR bqk -> z06
nrd XOR fgs -> wpb
frj XOR qhw -> z04
bqk OR frj -> z07
y03 OR x01 -> nrd
hwm AND bqk -> z03
tgd XOR rvg -> z12
tnw OR pbm -> gnj
`

var (
	reReg  = regexp.MustCompile(`^([a-z0-9]{3}): ([01])$`)
	reGate = regexp.MustCompile(`^([a-z0-9]{3}) ([ANDXOR]+) ([a-z0-9]{3}) -> ([a-z0-9]{3})`)
)

type Op byte

const (
	AND = Op(iota)
	OR
	XOR
)

type gate struct {
	a, b, o string
	op      Op
}

type wire struct {
	name  string
	value bool
}

func main() {
	run(sample1Gates, sample1Regs)
	run(sample2Gates, sample2Regs)
}

func parseGates(input string) []gate {
	var gates []gate
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if reGate.MatchString(line) {
			matches := reGate.FindStringSubmatch(line)
			switch matches[2] {
			case "AND":
				gates = append(gates, gate{a: matches[1], b: matches[3], o: matches[4], op: AND})
			case "OR":
				gates = append(gates, gate{a: matches[1], b: matches[3], o: matches[4], op: OR})
			case "XOR":
				gates = append(gates, gate{a: matches[1], b: matches[3], o: matches[4], op: XOR})
			default:
				panic(matches[2])
			}
		} else {
			panic(line)
		}
	}
	return gates
}

func parseRegs(input string) []wire {
	var wires []wire
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if reReg.MatchString(line) {
			matches := reReg.FindStringSubmatch(line)
			var value bool
			switch matches[2] {
			case "0":
				value = false
			case "1":
				value = true
			default:
				panic(matches[2])
			}
			wires = append(wires, wire{name: matches[1], value: value})
		} else {
			panic(line)
		}
	}
	return wires
}

func run(gatesStr string, regsStr string) {
	gates := parseGates(gatesStr)
	regs := parseRegs(regsStr)

	// map all gates by inputs
	inputToGate := map[string][]*gate{}
	for gi := range gates {
		g := &gates[gi]
		inputToGate[g.a] = append(inputToGate[g.a], g)
		inputToGate[g.b] = append(inputToGate[g.b], g)
	}

	gateInputs := map[string][]bool{}

	// iterative work queue to propagate computation
	for ri := 0; ri < len(regs); ri++ {
		reg := regs[ri]
		for _, g := range inputToGate[reg.name] {
			output := g.o
			switch len(gateInputs[output]) {
			case 0:
				gateInputs[output] = append(gateInputs[output], reg.value)
			case 1:
				gateInputs[output] = append(gateInputs[output], reg.value)
				var value bool
				switch g.op {
				case AND:
					value = gateInputs[output][0] && gateInputs[output][1]
				case OR:
					value = gateInputs[output][0] || gateInputs[output][1]
				case XOR:
					value = gateInputs[output][0] != gateInputs[output][1]
				default:
					panic(g.op)
				}
				regs = append(regs, wire{name: output, value: value})
			default:
				panic("too many inputs")
			}
		}
	}

	// every wire should now be resolved
	slices.SortFunc(regs, func(a, b wire) int {
		return strings.Compare(b.name, a.name) // sort from z99 down
	})

	var output int
	for _, w := range regs {
		if w.name[0] == 'z' {
			output <<= 1
			if w.value {
				output |= 1
			}
		}
	}
	fmt.Println(output)
}
