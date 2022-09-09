package kademlia

import (
	"math/big"
	"testing"
)

type testContact struct {
	id *big.Int
}

func (c *testContact) ID() *big.Int {
	return c.id
}

func mustNewTestContactFromString(s string) Contact {
	id, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic("invalid contact id")
	}
	return &testContact{id: id}
}

func TestEmptyBucketDepth(t *testing.T) {
	const (
		B = 3
		k = 1
	)
	b := newBucket(new(big.Int), bitMask(B))

	if b.depth() != 0 {
		t.Fail()
	}
}

func TestBucketDepth(t *testing.T) {
	const (
		B = 17
		k = 7
	)
	b := newBucket(new(big.Int), bitMask(B))

	for _, test := range []struct {
		input    string
		expected int
	}{
		{expected: 16, input: "3faa"}, // 00111111 10101010
		{expected: 5, input: "3855"},  // 00111000 01010101
		{expected: 4, input: "30aa"},  // 00110000 10101010
		{expected: 3, input: "2055"},  // 00100000 01010101
		{expected: 2, input: "00aa"},  // 00000000 10101010
		{expected: 1, input: "4055"},  // 01000000 01010101
		{expected: 0, input: "80aa"},  // 10000000 10101010
	} {
		b.put(mustNewTestContactFromString(test.input), k)
		if d := b.depth(); d != test.expected {
			t.Errorf("expected %d, got %d, %s", test.expected, d, test.input)
		}
	}
}
