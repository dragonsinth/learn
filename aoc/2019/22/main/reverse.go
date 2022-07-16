package main

import "fmt"

func newReverseDeck(ops []op, count int64) *reverseDeck {
	rd := &reverseDeck{
		count: count,
	}
	for i := int64(0); i < 100; i++ {
		rd.inverses[i] = modInverse(i, count)
	}
	r := make([]rop, len(ops))
	for i, in := range ops {
		arg := int64(in.arg)
		switch in.code {
		case 1:
			arg = rd.mod(count - arg)
		case 2:
			arg = rd.inverses[arg]
		}
		r[len(r)-i-1] = rop{
			code: in.code,
			arg:  arg,
		}
	}
	rd.rops = r
	rd.N, rd.K = rd.processProgram()
	fmt.Println(rd.N, rd.K)
	return rd
}

type reverseDeck struct {
	count    int64
	inverses [100]int64
	rops     []rop
	N, K     int64
}

func (rd *reverseDeck) Solve(slot int64, n int64) int64 {
	// Xn++ = Xn * N + K

	// Solve X*N^n
	lhs := rd.modMul(slot, rd.modPow(rd.N, n))

	// Solve K(1 + N + N^2 ... N^(n-1)
	// K(1 - N^n)/(1 - N)
	denom := rd.mod(1 - rd.N)
	denomInv := modInverse(denom, rd.count)
	if denomInv == 0 {
		panic("here")
	}
	num := rd.mod(1 - rd.modPow(rd.N, n))
	rhs := rd.modMul(rd.K, rd.modMul(num, denomInv))

	return (lhs + rhs) % rd.count
}

type rop struct {
	code int
	arg  int64
}

func (rd *reverseDeck) processProgram() (int64, int64) {
	N, K := int64(1), int64(0)
	for _, o := range rd.rops {
		switch o.code {
		case 0:
			N = rd.mod(-N)
			K = rd.mod(-K)
			K--
		case 1:
			K = rd.mod(K - o.arg)
		case 2:
			N = rd.modMul(N, o.arg)
			K = rd.modMul(K, o.arg)
		}
	}

	return N, K
}

func (rd *reverseDeck) reverse(in int64) int64 {
	return rd.mod(rd.modMul(in, rd.N) + rd.K)
}

func (rd *reverseDeck) reverseOld(slot int64) int64 {
	for _, o := range rd.rops {
		switch o.code {
		case 0:
			slot = rd.reverseDeal(slot)
		case 1:
			slot = rd.reverseCut(slot, o.arg)
		case 2:
			slot = rd.reverseIncrement(slot, o.arg)
		}

		if slot < 0 {
			panic(slot)
		}
	}
	return slot
}

func (rd *reverseDeck) reverseDeal(slot int64) int64 {
	// e.g. 9 -> 0, 8 -> 1
	// e.g  0 -> 9, 1 -> 8
	return rd.mod(-1 - slot)
}

func (rd *reverseDeck) reverseCut(slot int64, v int64) int64 {
	return rd.mod(slot - v)
}

func (rd *reverseDeck) reverseIncrement(slot int64, v int64) int64 {
	return rd.modMul(slot, v)
}

func (rd *reverseDeck) mod(a int64) int64 {
	m := rd.count
	return (a%m + m) % m
}

func (rd *reverseDeck) modMul(a, b int64) int64 {
	m := rd.count
	res := int64(0)
	a = a % m
	for b > 0 {
		// If b is odd, add 'a' to result
		if b%2 == 1 {
			res = (res + a) % m
		}

		// Multiply 'a' with 2
		a = (a * 2) % m

		// Divide b by 2
		b /= 2
	}
	return res % m
}

func (rd *reverseDeck) modPow(a, e int64) int64 {
	r := int64(1)
	for e > 0 {
		// If e is odd, multiply a with result
		if e&1 != 0 {
			r = rd.modMul(r, a)
		}
		// e must be even now
		e = e >> 1
		a = rd.modMul(a, a)
	}
	return r
}

// Function to find modulo inverse of a
func modInverse(a, m int64) int64 {
	g, x, _ := gcdExtended(a, m)
	if g != 1 {
		return 0
	}

	// m is added to handle negative x
	return (x%m + m) % m
}

// Function for extended Euclidean Algorithm
func gcdExtended(a, b int64) (int64, int64, int64) {
	if a == 0 {
		return b, 0, 1
	}
	gcd, x, y := gcdExtended(b%a, a)
	// Update x and y using results of recursive
	return gcd, y - (b/a)*x, x
}
