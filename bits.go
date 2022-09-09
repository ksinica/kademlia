package kademlia

import (
	"math/big"
	"unsafe"
)

func bitAlign(n, bits int) int {
	return n + ((bits - (n % bits)) % bits)
}

func bitMask(count int) *big.Int {
	const (
		W = int(unsafe.Sizeof(big.Word(0))) << 3
		B = (1 << W) - 1
	)
	q, r := count/W, count%W

	bits := make([]big.Word, 0, q)
	for i := 0; i < q; i++ {
		bits = append(bits, B)
	}

	if r > 0 {
		var x big.Word
		for i := 0; i < r; i++ {
			x |= 1 << i
		}
		bits = append(bits, x)
	}

	return new(big.Int).SetBits(bits)
}
