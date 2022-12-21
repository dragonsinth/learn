package main

import (
	"fmt"
	"io"
	"strconv"
)

type expr interface {
	Tok() string
	Eval() int
	Simplify() expr
	Code(out io.Writer)
	Reverse(Y int) int
}

type literal struct {
	key string
	val int
}

func (l literal) Tok() string {
	return l.key
}

func (l literal) Eval() int {
	return l.val
}

func (l literal) Simplify() expr {
	if l.key == "humn" {
		return l
	}
	return simpleLit(l.val)
}

func (l literal) Code(out io.Writer) {
	_, _ = fmt.Fprintf(out, "%s := %d\n", l.key, l.val)
}

func (l literal) Reverse(Y int) int {
	if l.key == "humn" {
		return Y
	}
	panic("here")
}

var _ expr = literal{}

type simpleLit int

func (l simpleLit) Tok() string {
	return strconv.Itoa(int(l))
}

func (l simpleLit) Eval() int {
	return int(l)
}

func (l simpleLit) Simplify() expr {
	return l
}

func (l simpleLit) Code(out io.Writer) {
}

func (l simpleLit) Reverse(_ int) int {
	panic("here")
}

var _ expr = simpleLit(0)

type operation struct {
	key string
	lhs expr
	op  byte
	rhs expr
}

func (o operation) Tok() string {
	return o.key
}

func (o operation) Eval() int {
	switch o.op {
	case '+':
		return o.lhs.Eval() + o.rhs.Eval()
	case '-':
		return o.lhs.Eval() - o.rhs.Eval()
	case '*':
		return o.lhs.Eval() * o.rhs.Eval()
	case '/':
		return o.lhs.Eval() / o.rhs.Eval()
	case '=':
		if o.lhs.Eval() == o.rhs.Eval() {
			return 1
		}
		return 0
	default:
		panic(o.op)
	}
}

func (o operation) Simplify() expr {
	lhs := o.lhs.Simplify()
	rhs := o.rhs.Simplify()
	_, lhsLit := lhs.(simpleLit)
	_, rhsLit := rhs.(simpleLit)

	ret := operation{
		key: o.key,
		lhs: lhs,
		op:  o.op,
		rhs: rhs,
	}
	if lhsLit && rhsLit {
		return simpleLit(ret.Eval())
	}
	return ret
}

func (o operation) Code(out io.Writer) {
	o.lhs.Code(out)
	o.rhs.Code(out)
	_, _ = fmt.Fprintf(out, "%s := %s %s %s\n", o.key, o.lhs.Tok(), string(o.op), o.rhs.Tok())
}

func (o operation) Reverse(Y int) int {
	// one side better be a simple literal
	var X int
	var exp expr
	_, lhsLit := o.lhs.(simpleLit)
	if lhsLit {
		X, exp = int(o.lhs.(simpleLit)), o.rhs
	} else {
		X, exp = int(o.rhs.(simpleLit)), o.lhs
	}

	switch o.op {
	case '+':
		// commutative
		// Y = X + exp -> Y - X = exp
		return exp.Reverse(Y - X)
	case '-':
		// non-commutative
		if lhsLit {
			// Y = X - exp -> exp = X - Y
			return exp.Reverse(X - Y)
		} else {
			// Y = exp - X -> Y + X = exp
			return exp.Reverse(Y + X)
		}
	case '*':
		// commutative
		// Y = X * exp -> Y / X = exp
		return exp.Reverse(Y / X)
	case '/':
		// non-commutative
		if lhsLit {
			// Y = X / exp -> exp = X / Y
			return exp.Reverse(X / Y)
		} else {
			// Y = exp / X -> Y * X = exp
			return exp.Reverse(Y * X)
		}
	case '=':
		if Y != 1 {
			panic(Y)
		}
		// true = exp == X
		return exp.Reverse(X)
	default:
		panic(o.op)
	}
}

var _ expr = operation{}
