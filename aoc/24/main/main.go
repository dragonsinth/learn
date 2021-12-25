package main

import (
	"fmt"
)

func main() {
	results := findAll([]int64{}, 0, 0, 0, []int{4, 2, 6, 2}, steps)
	for _, r := range results {
		fmt.Printf("%+v = ", r.path)

		// See what the first 4 steps yields
		var x, y, z int64
		for i := 0; i < 4; i++ {
			_, x, y, z = steps[i](r.path[i], x, y, z)
		}
		fmt.Println(z, base26(z))
	}
	fmt.Println(len(results))
}

func findAll(path []int64, x, y, z int64, parts []int, fs []computefunc) []result {
	var ret []result

	results := findSegments(path, x, y, z, fs[0:parts[0]])
	if len(parts) > 1 {
		// For each result, recurse the next.
		for _, r := range results {
			rec := findAll(r.path, 0, 0, r.z, parts[1:], fs[parts[0]:])
			ret = append(ret, rec...)
		}
	} else {
		// Merge the results
		ret = append(ret, results...)
	}

	return ret
}

func findSegments(path []int64, x, y, z int64, fs []computefunc) []result {
	if len(fs) == 0 {
		keep := x == 0 && y == 0 && z < 1000
		if keep && len(path) == 14 {
			keep = z == 0
		}
		if keep {
			return []result{
				{
					z:    z,
					path: append([]int64{}, path...),
				},
			}
		}
		return nil
	}
	var ret []result
	for i := 1; i <= 9; i++ {
		_, x, y, z := fs[0](int64(i), x, y, z)
		rec := findSegments(append(path, int64(i)), x, y, z, fs[1:])
		ret = append(ret, rec...)
	}
	return ret
}

type result struct {
	z    int64
	path []int64
}

func base26(z int64) string {
	var buf []byte
	for z > 0 {
		b := z % 26
		z = z / 26
		buf = append([]byte{byte(b + 'a')}, buf...)
	}
	return string(buf)
}
