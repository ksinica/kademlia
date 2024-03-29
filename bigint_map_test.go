package kademlia

import (
	"math/big"
	"math/rand"
	"testing"
)

func newRand() *rand.Rand {
	return rand.New(rand.NewSource(0))
}

func randBytes(rand *rand.Rand, size int) []byte {
	b := make([]byte, size)
	for i := range b {
		b[i] = byte(rand.Intn(0xff))
	}
	return b
}

func randBigInts(n, size int) (ret []*big.Int) {
	m, rand := make(map[string]struct{}), newRand()
	for len(m) < n {
		m[string(randBytes(rand, size))] = struct{}{}
	}

	for k := range m {
		ret = append(ret, new(big.Int).SetBytes([]byte(k)))
	}
	return
}

func BenchmarkIntMapSet256(b *testing.B) {
	keys := randBigInts(b.N, 256/8)
	var m bigIntMap[struct{}]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.set(keys[i], struct{}{})
	}
}

func BenchmarkIntMapContains256(b *testing.B) {
	keys := randBigInts(b.N, 256/8)
	var m bigIntMap[struct{}]

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
