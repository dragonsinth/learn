package main

import (
	"testing"
)

func TestModInverse(t *testing.T) {
	m := modInverse(39, 10007)
	t.Log(m)
	for i := int64(0); i < 10007; i++ {
		v := (i * 39) % 10007
		rv := (v * m) % 10007
		if i != rv {
			t.Errorf("%d != %d", i, rv)
		}
	}

	t.Log(gcdExtended(102, 38))
	t.Log(gcdExtended(42823, 6409))
}
