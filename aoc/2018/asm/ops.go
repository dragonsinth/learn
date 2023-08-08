package asm

import "fmt"

var Instrs = map[string]OpFun{
	"addi": addi,
	"addr": addr,
	"bani": bani,
	"banr": banr,
	"bori": bori,
	"borr": borr,
	"eqir": eqir,
	"eqri": eqri,
	"eqrr": eqrr,
	"gtir": gtir,
	"gtri": gtri,
	"gtrr": gtrr,
	"muli": muli,
	"mulr": mulr,
	"seti": seti,
	"setr": setr,
}

func (in Instr) ToString(names [6]string) string {
	dst := names[in.Ops[2]]

	switch in.Name {
	case "addi":
		return fmt.Sprint(dst, " = ", names[in.Ops[0]], " + ", in.Ops[1])
	case "addr":
		return fmt.Sprint(dst, " = ", names[in.Ops[0]], " + ", names[in.Ops[1]])
	case "bani":
		return fmt.Sprint(dst, " = ", names[in.Ops[0]], " & ", in.Ops[1])
	case "banr":
		return fmt.Sprint(dst, " = ", names[in.Ops[0]], " & ", names[in.Ops[1]])
	case "bori":
		return fmt.Sprint(dst, " = ", names[in.Ops[0]], " | ", in.Ops[1])
	case "borr":
		return fmt.Sprint(dst, " = ", names[in.Ops[0]], " | ", names[in.Ops[1]])
	case "eqir":
		return fmt.Sprint(dst, " if ", in.Ops[0], " == ", names[in.Ops[1]])
	case "eqri":
		return fmt.Sprint(dst, " if ", names[in.Ops[0]], " == ", in.Ops[1])
	case "eqrr":
		return fmt.Sprint(dst, " if ", names[in.Ops[0]], " == ", names[in.Ops[1]])
	case "gtir":
		return fmt.Sprint(dst, " if ", in.Ops[0], " > ", names[in.Ops[1]])
	case "gtri":
		return fmt.Sprint(dst, " if ", names[in.Ops[0]], " > ", in.Ops[1])
	case "gtrr":
		return fmt.Sprint(dst, " if ", names[in.Ops[0]], " > ", names[in.Ops[1]])
	case "muli":
		return fmt.Sprint(dst, " = ", names[in.Ops[0]], " * ", in.Ops[1])
	case "mulr":
		return fmt.Sprint(dst, " = ", names[in.Ops[0]], " * ", names[in.Ops[1]])
	case "seti":
		return fmt.Sprint(dst, " = ", in.Ops[0])
	case "setr":
		return fmt.Sprint(dst, " = ", names[in.Ops[0]])
	default:
		panic(in.Name)
	}
}

func addr(ops [3]int, reg Registers) Registers {
	out := reg
	A := reg[ops[0]]
	B := reg[ops[1]]
	C := ops[2]
	out[C] = A + B
	return out
}

func addi(ops [3]int, reg Registers) Registers {
	out := reg
	A := reg[ops[0]]
	B := ops[1]
	C := ops[2]
	out[C] = A + B
	return out
}

func mulr(ops [3]int, reg Registers) Registers {
	out := reg
	A := reg[ops[0]]
	B := reg[ops[1]]
	C := ops[2]
	out[C] = A * B
	return out
}

func muli(ops [3]int, reg Registers) Registers {
	out := reg
	A := reg[ops[0]]
	B := ops[1]
	C := ops[2]
	out[C] = A * B
	return out
}

func banr(ops [3]int, reg Registers) Registers {
	out := reg
	A := reg[ops[0]]
	B := reg[ops[1]]
	C := ops[2]
	out[C] = A & B
	return out
}

func bani(ops [3]int, reg Registers) Registers {
	out := reg
	A := reg[ops[0]]
	B := ops[1]
	C := ops[2]
	out[C] = A & B
	return out
}

func borr(ops [3]int, reg Registers) Registers {
	out := reg
	A := reg[ops[0]]
	B := reg[ops[1]]
	C := ops[2]
	out[C] = A | B
	return out
}

func bori(ops [3]int, reg Registers) Registers {
	out := reg
	A := reg[ops[0]]
	B := ops[1]
	C := ops[2]
	out[C] = A | B
	return out
}

func setr(ops [3]int, reg Registers) Registers {
	out := reg
	A := reg[ops[0]]
	C := ops[2]
	out[C] = A
	return out
}

func seti(ops [3]int, reg Registers) Registers {
	out := reg
	A := ops[0]
	C := ops[2]
	out[C] = A
	return out
}

func gtir(ops [3]int, reg Registers) Registers {
	out := reg
	A := ops[0]
	B := reg[ops[1]]
	C := ops[2]
	out[C] = toInt(A > B)
	return out
}

func gtri(ops [3]int, reg Registers) Registers {
	out := reg
	A := reg[ops[0]]
	B := ops[1]
	C := ops[2]
	out[C] = toInt(A > B)
	return out
}

func gtrr(ops [3]int, reg Registers) Registers {
	out := reg
	A := reg[ops[0]]
	B := reg[ops[1]]
	C := ops[2]
	out[C] = toInt(A > B)
	return out
}

func eqir(ops [3]int, reg Registers) Registers {
	out := reg
	A := ops[0]
	B := reg[ops[1]]
	C := ops[2]
	out[C] = toInt(A == B)
	return out
}

func eqri(ops [3]int, reg Registers) Registers {
	out := reg
	A := reg[ops[0]]
	B := ops[1]
	C := ops[2]
	out[C] = toInt(A == B)
	return out
}

func eqrr(ops [3]int, reg Registers) Registers {
	out := reg
	A := reg[ops[0]]
	B := reg[ops[1]]
	C := ops[2]
	out[C] = toInt(A == B)
	return out
}

func toInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
