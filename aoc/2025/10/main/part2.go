package main

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
)

func (p *puzzle) Part2() int {
	buttonCount := len(p.buttons)

	// convert puzzle into a matrix
	// we need a row for each joltage plus a final row for summation
	var rows [][]int
	for joltageIt := 0; joltageIt < len(p.joltages); joltageIt++ {
		var row []int
		for buttonIt := 0; buttonIt < buttonCount; buttonIt++ {
			// Does this button increase the current joltage?
			if slices.Contains(p.buttons[buttonIt], joltageIt) {
				row = append(row, 1)
			} else {
				row = append(row, 0)
			}
		}
		// Add the negative joltage on the RHS as a constant
		row = append(row, -int(p.joltages[joltageIt]))
		row = append(row, 0) // factor of X
		rows = append(rows, row)
	}

	// final row: each button count sums to 0 + 1K
	rows = append(rows, append(slices.Repeat([]int{1}, buttonCount), 0, -1))

	// (rows)
	doGaussian(rows, buttonCount)

	success := false
	defer func() {
		if !success {
			printMatrix(rows)
		}
	}()
	ret, success := trySolveXDirect(rows, buttonCount)
	if !success {
		ret, success = minXGeneral(rows, buttonCount)
	}
	if ret <= 0 {
		success = false
	}
	return ret
}

func printMatrix(rows [][]int) {
	for _, row := range rows {
		out, _ := json.Marshal(row)
		fmt.Println(string(out))
	}
	fmt.Println()
}

func goprintMatrix(rows [][]int) {
	fmt.Printf("%#v", rows)
	fmt.Println()
}

// -------------------------------------------------------------- //
// BENEATH THIS LINE YOU WILL FIND YOURSELF IN THE DEPTHS OF
// MASSIVE AI CODING ASSISTANCE. I LEARNED SOME THINGS ABOUT
// SYSTEMS OF LINEAR EQUATIONS BUT I BARELY KNOW WHAT I'M DOING.
// -------------------------------------------------------------- //

func doGaussian(rows [][]int, N int) {
	if len(rows) == 0 {
		return
	}
	m := len(rows)
	n := len(rows[0])

	h, k := 0, 0
	for h < m && k < N {
		// Find the pivot row in column k starting at row h
		imax := h
		best := abs(rows[h][k])
		for i := h + 1; i < m; i++ {
			v := abs(rows[i][k])
			if v > best {
				best = v
				imax = i
			}
		}
		// If the pivot is zero, move to next column
		if rows[imax][k] == 0 {
			k++
			continue
		}
		// Swap current row h with pivot row imax
		if imax != h {
			rows[h], rows[imax] = rows[imax], rows[h]
		}

		// Normalize the pivot row so that pivot becomes 1
		pivot := rows[h][k]
		if pivot < 0 {
			for j := k; j < n; j++ {
				rows[h][j] = -rows[h][j]
			}
			pivot = -pivot
		}

		// Eliminate this column in all other rows
		for i := 0; i < m; i++ {
			if i == h {
				continue
			}
			f := rows[i][k]
			if f == 0 {
				continue
			}

			// rows[i] = rows[i]*p - rows[h]*f
			for j := 0; j < n; j++ {
				rows[i][j] = rows[i][j]*pivot - f*rows[h][j]
			}

			// reduce row i by gcd to control growth
			g := 0
			for j := 0; j < n; j++ {
				g = gcd(g, abs(rows[i][j]))
			}
			if g > 1 {
				for j := 0; j < n; j++ {
					rows[i][j] /= g
				}
			}
		}
		h++
		k++
	}
}

// Try to read a direct equality for X from any row with all-zero variable cols.
func trySolveXDirect(rows [][]int, N int) (x int, ok bool) {
	zeroes := make([]int, N)
	xFound := false
	xValue := 0
	for r := 0; r < len(rows); r++ {
		if slices.Equal(rows[r][:N], zeroes) {
			constCol := rows[r][N]
			xFac := rows[r][N+1]
			if xFac == 0 {
				continue
			}
			num := -constCol
			den := xFac
			// exact division check; adjust if you allow rationals
			if den == 0 {
				continue
			}
			if num%den != 0 {
				return 0, false
			} // inconsistent with integer model
			xVal := num / den
			if !xFound {
				xFound = true
				xValue = xVal
			} else if xValue != xVal {
				// conflicting X equations
				panic(xValue)
			}
		}
	}
	return xValue, xFound
}

// Build mapping from pivot column -> row index by finding the first nonzero
// among the first N columns for each row. This assumes the matrix is in a
// row-echelon-like form after doGaussian, so each pivot row has a unique
// leading column within [0..N).
func detectPivots(rows [][]int, N int) map[int]int {
	piv := make(map[int]int)
	for r := 0; r < len(rows); r++ {
		for c := 0; c < N; c++ {
			if rows[r][c] != 0 {
				piv[c] = r
				break
			}
		}
	}
	return piv
}

// Represents α·X + β1·t1 + β2·t2 + γ ≥ 0
type aff struct {
	alpha int
	beta  []int
	gamma int

	xFactor  int   // rows[r][N+1]
	constCol int   // rows[r][N]
	free     []int // rows[r][frees[j]]
	mod      int   // abs(ci) where ci = rows[r][i]; for free vars set to 1
}

func buildForms(rows [][]int, N int) ([]aff, []int) {
	piv := detectPivots(rows, N)

	var frees []int
	for i := 0; i < N; i++ {
		if _, ok := piv[i]; !ok {
			frees = append(frees, i)
		}
	}

	forms := make([]aff, N)
outer:
	for i := 0; i < N; i++ {
		form := aff{
			alpha: 0,
			beta:  make([]int, len(frees)),
			gamma: 0,

			xFactor:  0,
			constCol: 0,
			free:     make([]int, len(frees)),
			mod:      1,
		}
		// Free var itself → vi = ti ≥ 0 (we’ll enforce t>=0 as loop bounds)
		for j := 0; j < len(frees); j++ {
			if i == frees[j] {
				form.beta[j] = 1
				forms[i] = form
				continue outer
			}
		}
		r, ok := piv[i]
		if !ok {
			panic(fmt.Errorf("non-pivot col %d but not listed as free", i))
		}
		ci := rows[r][i]
		form.mod = abs(ci)
		form.constCol = rows[r][N]
		form.xFactor = rows[r][N+1]

		s := 1
		if ci > 0 {
			s = -1
		}
		form.alpha = s * form.xFactor
		for j := 0; j < len(frees); j++ {
			freeCoeff := rows[r][frees[j]]
			form.beta[j] = s * freeCoeff
			form.free[j] = freeCoeff
		}
		form.gamma = s * form.constCol
		forms[i] = form
	}
	return forms, frees
}

func minXGeneral(rows [][]int, N int) (int, bool) {
	forms, frees := buildForms(rows, N)
	tFree := make([]int, len(frees))
	return recurseMinX(forms, N, tFree, 0)
}

func recurseMinX(forms []aff, N int, tFree []int, depth int) (int, bool) {
	if depth < len(tFree) {
		ret := math.MaxInt
		retOk := false
		for t := 0; t < 500; t++ {
			tFree[depth] = t
			if val, ok := recurseMinX(forms, N, tFree, depth+1); ok {
				ret = min(ret, val)
				retOk = true
			}
		}
		return ret, retOk
	}

	// All variables are pinned, evaluate
	lb := math.MinInt
	ub := math.MaxInt
	for i := 0; i < N; i++ {
		a, b, g := forms[i].alpha, forms[i].beta, forms[i].gamma
		// rhs = -b*t - g
		rhs := -g
		for i := range b {
			rhs -= b[i] * tFree[i]
		}
		if a > 0 {
			lb = max(lb, divCeil(rhs, a))
		}
		if a < 0 {
			ub = min(ub, divFloor(rhs, a))
		}
		if a == 0 {
			// Constraint reduces to b*t + g >= 0
			// Therefore rhs (-b*t - g) <= 0
			if rhs > 0 {
				// infeasible for this tFree
				lb = 1
				ub = 0
			}
		}
	}
	if lb <= ub {
		// try smallest X respecting all congruences
		for X := lb; X <= ub; X++ {
			ok := true
			for i := 0; i < N && ok; i++ {
				m := forms[i].mod
				if m <= 1 {
					continue
				} // free var or no restriction
				v := forms[i].xFactor*X + forms[i].constCol
				for j := range forms[i].free {
					v += forms[i].free[j] * tFree[j]
				}
				if v%m != 0 {
					ok = false
					break
				}
			}
			if ok {
				return X, true
			}
		}
	}
	return 0, false
}
