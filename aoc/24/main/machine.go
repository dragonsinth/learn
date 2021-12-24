package main

type computefunc func(w, x, y, z int64) (int64, int64, int64, int64)

var steps = []computefunc{
	func(w, x, y, z int64) (int64, int64, int64, int64) {
		x = x * 0
		x = x + z
		x = x % 26
		z = z / 1
		x = x + 12
		x = eq(x, w)
		x = eq(x, 0)
		y = y * 0
		y = y + 25
		y = y * x
		y = y + 1
		z = z * y
		y = y * 0
		y = y + w
		y = y + 6
		y = y * x
		z = z + y
		return w, x, y, z
	},
	func(w, x, y, z int64) (int64, int64, int64, int64) {
		x = x * 0
		x = x + z
		x = x % 26
		z = z / 1
		x = x + 10
		x = eq(x, w)
		x = eq(x, 0)
		y = y * 0
		y = y + 25
		y = y * x
		y = y + 1
		z = z * y
		y = y * 0
		y = y + w
		y = y + 6
		y = y * x
		z = z + y
		return w, x, y, z
	},
	func(w, x, y, z int64) (int64, int64, int64, int64) {
		x = x * 0
		x = x + z
		x = x % 26
		z = z / 1
		x = x + 13
		x = eq(x, w)
		x = eq(x, 0)
		y = y * 0
		y = y + 25
		y = y * x
		y = y + 1
		z = z * y
		y = y * 0
		y = y + w
		y = y + 3
		y = y * x
		z = z + y
		return w, x, y, z
	},
	func(w, x, y, z int64) (int64, int64, int64, int64) {
		x = x * 0
		x = x + z
		x = x % 26
		z = z / 26
		x = x + -11
		x = eq(x, w)
		x = eq(x, 0)
		y = y * 0
		y = y + 25
		y = y * x
		y = y + 1
		z = z * y
		y = y * 0
		y = y + w
		y = y + 11
		y = y * x
		z = z + y
		return w, x, y, z
	},
	func(w, x, y, z int64) (int64, int64, int64, int64) {
		x = x * 0
		x = x + z
		x = x % 26
		z = z / 1
		x = x + 13
		x = eq(x, w)
		x = eq(x, 0)
		y = y * 0
		y = y + 25
		y = y * x
		y = y + 1
		z = z * y
		y = y * 0
		y = y + w
		y = y + 9
		y = y * x
		z = z + y
		return w, x, y, z
	},
	func(w, x, y, z int64) (int64, int64, int64, int64) {
		x = x * 0
		x = x + z
		x = x % 26
		z = z / 26
		x = x + -1
		x = eq(x, w)
		x = eq(x, 0)
		y = y * 0
		y = y + 25
		y = y * x
		y = y + 1
		z = z * y
		y = y * 0
		y = y + w
		y = y + 3
		y = y * x
		z = z + y
		return w, x, y, z
	},
	func(w, x, y, z int64) (int64, int64, int64, int64) {
		x = x * 0
		x = x + z
		x = x % 26
		z = z / 1
		x = x + 10
		x = eq(x, w)
		x = eq(x, 0)
		y = y * 0
		y = y + 25
		y = y * x
		y = y + 1
		z = z * y
		y = y * 0
		y = y + w
		y = y + 13
		y = y * x
		z = z + y
		return w, x, y, z
	},
	func(w, x, y, z int64) (int64, int64, int64, int64) {
		x = x * 0
		x = x + z
		x = x % 26
		z = z / 1
		x = x + 11
		x = eq(x, w)
		x = eq(x, 0)
		y = y * 0
		y = y + 25
		y = y * x
		y = y + 1
		z = z * y
		y = y * 0
		y = y + w
		y = y + 6
		y = y * x
		z = z + y
		return w, x, y, z
	},
	func(w, x, y, z int64) (int64, int64, int64, int64) {
		x = x * 0
		x = x + z
		x = x % 26
		z = z / 26
		x = x + 0
		x = eq(x, w)
		x = eq(x, 0)
		y = y * 0
		y = y + 25
		y = y * x
		y = y + 1
		z = z * y
		y = y * 0
		y = y + w
		y = y + 14
		y = y * x
		z = z + y
		return w, x, y, z
	},
	func(w, x, y, z int64) (int64, int64, int64, int64) {
		x = x * 0
		x = x + z
		x = x % 26
		z = z / 1
		x = x + 10
		x = eq(x, w)
		x = eq(x, 0)
		y = y * 0
		y = y + 25
		y = y * x
		y = y + 1
		z = z * y
		y = y * 0
		y = y + w
		y = y + 10
		y = y * x
		z = z + y
		return w, x, y, z
	},
	func(w, x, y, z int64) (int64, int64, int64, int64) {
		x = x * 0
		x = x + z
		x = x % 26
		z = z / 26
		x = x + -5
		x = eq(x, w)
		x = eq(x, 0)
		y = y * 0
		y = y + 25
		y = y * x
		y = y + 1
		z = z * y
		y = y * 0
		y = y + w
		y = y + 12
		y = y * x
		z = z + y
		return w, x, y, z
	},
	func(w, x, y, z int64) (int64, int64, int64, int64) {
		x = x * 0
		x = x + z
		x = x % 26
		z = z / 26
		x = x + -16
		x = eq(x, w)
		x = eq(x, 0)
		y = y * 0
		y = y + 25
		y = y * x
		y = y + 1
		z = z * y
		y = y * 0
		y = y + w
		y = y + 10
		y = y * x
		z = z + y
		return w, x, y, z
	},
	func(w, x, y, z int64) (int64, int64, int64, int64) {
		x = x * 0
		x = x + z
		x = x % 26
		z = z / 26
		x = x + -7
		x = eq(x, w)
		x = eq(x, 0)
		y = y * 0
		y = y + 25
		y = y * x
		y = y + 1
		z = z * y
		y = y * 0
		y = y + w
		y = y + 11
		y = y * x
		z = z + y
		return w, x, y, z
	},
	func(w, x, y, z int64) (int64, int64, int64, int64) {
		x = x * 0
		x = x + z
		x = x % 26
		z = z / 26
		x = x + -11
		x = eq(x, w)
		x = eq(x, 0)
		y = y * 0
		y = y + 25
		y = y * x
		y = y + 1
		z = z * y
		y = y * 0
		y = y + w
		y = y + 15
		y = y * x
		z = z + y
		return w, x, y, z
	},
}

func eq(a int64, b int64) int64 {
	if a == b {
		return 1
	}
	return 0
}
