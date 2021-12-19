package main

import (
	"fmt"
	"regexp"
	"sort"
)

var input = sample

var (
	scannerLine = regexp.MustCompile(`^--- scanner (\d+) ---$`)
)

type Dim int

const (
	X = Dim(0)
	Y = Dim(1)
	Z = Dim(2)
)

type scanner struct {
	id        int
	points    []point
	offset    *point
	absPoints []point
}

func (s *scanner) adjust(abs point) {
	s.offset = &abs
	for _, p := range s.points {
		s.absPoints = append(s.absPoints, point{
			p[0] + abs[X],
			p[1] + abs[Y],
			p[2] + abs[Z],
		})
	}
}

func (s scanner) String() string {
	if s.offset != nil {
		return fmt.Sprintf("%d:%s", s.id, s.offset)
	} else {
		return fmt.Sprintf("%d:???", s.id)
	}
}

type point [3]int

func (p point) String() string {
	return fmt.Sprintf("%d,%d,%d", p[X], p[Y], p[Z])
}

func main() {
	scanners := parseScanners(input)

	scanners[0].adjust(point{})
	aligned := []*scanner{scanners[0]}
	unalighed := scanners[1:]

	fmt.Println(aligned[0])
	fmt.Println(unalighed[0])

	for len(unalighed) > 0 {
		// Find an unaligned scanner we can move.
	outer:
		for uit := range unalighed {
			for ait := range aligned {
				if rotated := tryAlign(12, aligned[ait], unalighed[uit]); rotated != nil {
					fmt.Printf("aligned %s to %s\n", rotated, aligned[ait])
					aligned = append(aligned, rotated)
					unalighed = append(unalighed[:uit], unalighed[uit+1:]...)
					break outer
				}

			}
		}
	}

	uniq := map[point]bool{}
	for _, sc := range aligned {
		for _, p := range sc.absPoints {
			uniq[p] = true
		}
	}
	fmt.Println(len(uniq))
}

func tryAlign(minPoints int, dst *scanner, src *scanner) *scanner {
	// Compute all permutations of the src scanner
	perms := permutations(src.points)
	for _, perm := range perms {
		if p := alignPoints(minPoints, dst.points, perm); p != nil {
			ret := &scanner{id: src.id, points: perm}
			// Must adjust p relative to the destination scanner
			p[X] += dst.offset[X]
			p[Y] += dst.offset[Y]
			p[Z] += dst.offset[Z]
			ret.adjust(*p)
			return ret
		}
	}

	return nil
}

func alignPoints(minPoints int, dst []point, src []point) *point {
	for dx := -2000; dx < 2000; dx++ {
		// see if we can get at least 12 X coords to align.
		if !alignPointsDim(X, minPoints, dst, src, dx) {
			continue
		}
		//fmt.Println("align X", dx)

		for dy := -2000; dy < 2000; dy++ {
			if !alignPointsDim(Y, minPoints, dst, src, dy) {
				continue
			}
			//fmt.Println("align Y")
			for dz := -2000; dz < 2000; dz++ {
				if !alignPointsDim(Z, minPoints, dst, src, dz) {
					continue
				}
				//fmt.Println("align Z")

				dp := point{dx, dy, dz}
				if alignPointsExact(minPoints, dst, src, dp) {
					//fmt.Println("align exact")
					return &dp
				} else {
					fmt.Println("miss")
					continue
				}
			}
		}
	}

	return nil
}

func alignPointsDim(dim Dim, minPoints int, dstPts []point, srcPts []point, deltaSrc int) bool {
	dst := make([]int, len(dstPts))
	for i, p := range dstPts {
		dst[i] = p[dim]
	}

	src := make([]int, len(srcPts))
	for i, p := range srcPts {
		src[i] = p[dim] + deltaSrc
	}

	sort.Ints(dst)
	sort.Ints(src)

	hits := 0
	for d, s := 0, 0; d < len(dst) && s < len(src); {
		if dst[d] < src[s] {
			d++
		} else if src[s] < dst[d] {
			s++
		} else {
			hits++
			d++
			s++
		}
	}

	return hits >= minPoints
}

func alignPointsExact(minPoints int, dstPts []point, srcPts []point, dp point) bool {
	dst := map[point]bool{}
	for _, p := range dstPts {
		dst[p] = true
	}

	hits := 0
	for _, srcPt := range srcPts {
		if dst[point{
			X: srcPt[X] + dp[X],
			Y: srcPt[Y] + dp[Y],
			Z: srcPt[Z] + dp[Z],
		}] {
			hits++
		}
	}

	return hits >= minPoints
}
