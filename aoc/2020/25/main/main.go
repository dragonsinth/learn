package main

import "fmt"

var (
	sample = []int{5764801, 17807724}
	data   = []int{14012298, 74241}
	input  = data
)

const (
	keyRoot = 7
	modulus = 20201227
)

func main() {
	secrets := []int{0, 0}
	found := 0
	for i := 0; found < 2; i++ {
		v := transform2(i, keyRoot)
		//fmt.Println(i, v, transform2(i, keyRoot))
		if v == input[0] {
			//fmt.Println(0, i)
			secrets[0] = i
			found++
		}
		if v == input[1] {
			//fmt.Println(1, i)
			secrets[1] = i
			found++
		}
	}

	fmt.Println(secrets)

	fmt.Println(transform(secrets[0], input[1]))
	fmt.Println(transform(secrets[1], input[0]))
}

func transform(secret int, in int) int {
	v := 1
	for i := 0; i < secret; i++ {
		v = (v * in) % modulus
	}
	return v
}

func transform2(secret int, in int) int {
	ret := 1
	for secret > 0 {
		if secret&1 != 0 {
			ret = (ret * in) % modulus
		}
		in = (in * in) % modulus
		secret >>= 1
	}
	return ret
}
