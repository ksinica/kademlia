package kademlia

import (
	"math/big"
	"testing"
)

// goos: linux
// goarch: amd64
// pkg: github.com/ksinica/kademlia
// cpu: AMD Ryzen 5 5600X 6-Core Processor
// BenchmarkIntMap-12   27292102   56.75 ns/op   24 B/op   1 allocs/op
func BenchmarkIntMap(b *testing.B) {
	var m intMap[struct{}]
	var answer big.Int
	answer.SetString("4242424242424242424242424242424242424242424242424242", 10)
	for i := 0; i < b.N; i++ {
		m.set(&answer, struct{}{})
		m.get(&answer)
	}
}
