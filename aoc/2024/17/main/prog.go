package main

type Program struct {
	Ops []byte

	A  int
	B  int
	C  int
	IP int

	Out []byte
}

var opFuncs = []func(*Program, byte){
	0: (*Program).adv,
	1: (*Program).bxl,
	2: (*Program).bst,
	3: (*Program).jnz,
	4: (*Program).bxc,
	5: (*Program).out,
	6: (*Program).bdv,
	7: (*Program).cdv,
}

func (p *Program) Run() []byte {
	for p.IP = 0; p.IP < len(p.Ops); p.IP += 2 {
		code, op := p.Ops[p.IP], p.Ops[p.IP+1]
		opFuncs[code](p, op)
	}
	return p.Out
}

func (p *Program) adv(op byte) {
	p.A = p.div(op)
}

func (p *Program) bxl(op byte) {
	p.B = p.B ^ int(op)
}

func (p *Program) bst(op byte) {
	p.B = p.combo(op) % 8
}

func (p *Program) jnz(op byte) {
	if p.A != 0 {
		p.IP = int(op)
		p.IP -= 2
	}
}

func (p *Program) bxc(_ byte) {
	p.B = p.B ^ p.C
}

func (p *Program) out(op byte) {
	p.Out = append(p.Out, byte(p.combo(op)%8))
}

func (p *Program) bdv(op byte) {
	p.B = p.div(op)
}

func (p *Program) cdv(op byte) {
	p.C = p.div(op)
}

func (p *Program) div(op byte) int {
	return p.A / (1 << p.combo(op))
}

func (p *Program) combo(op byte) int {
	switch op {
	case 0, 1, 2, 3:
		return int(op)
	case 4:
		return p.A
	case 5:
		return p.B
	case 6:
		return p.C
	default:
		panic(op)
	}
}
