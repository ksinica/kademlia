package kademlia

import (
	"math/big"
	"math/rand"
	"testing"
)

func newRand() *rand.Rand {
	return rand.New(rand.NewSource(1))
}

func randBytes(size int) []byte {
	b, rand := make([]byte, size), newRand()
	for i := range b {
		b[i] = byte(rand.Intn(0xff))
	}
	return b
}

func randBigInts(n, size int) (ret []*big.Int) {
	m := make(map[string]struct{})
	for len(m) < n {
		m[string(randBytes(size))] = struct{}{}
	}

	for k := range m {
		ret = append(ret, new(big.Int).SetBytes([]byte(k)))
	}
	return
}

func BenchmarkIntMapSet256(b *testing.B) {
	keys := randBigInts(b.N, 256/8)
	var m intMap[struct{}]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.set(keys[i], struct{}{})
	}
}

func BenchmarkIntMapContains256(b *testing.B) {
	keys := randBigInts(b.N, 256/8)
	var m intMap[struct{}]

	for _, k := range keys {
		m.set(k, struct{}{})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !m.contains(keys[i]) {
			panic("sould never happen")
		}
	}
}
