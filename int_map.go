package kademlia

import (
	"math/big"
	"runtime"
)

type intMapEntry[T any] struct {
	key   *big.Int
	value T
}

// intMap is a map with big.Int keys. It offers a convenient and optimized way
// to relate some arbitrary types with a big.Int.
type intMap[T any] struct {
	m map[string]intMapEntry[T]
}

func (m *intMap[T]) set(key *big.Int, value T) {
	if m.m == nil {
		m.m = make(map[string]intMapEntry[T])
	}
	m.m[string(unsafeToBytes(key))] = intMapEntry[T]{
		key:   key,
		value: value,
	}
	runtime.KeepAlive(key)
}

func (m *intMap[T]) get(key *big.Int) (value T, ok bool) {
	if m.m != nil {
		e, ok := m.m[string(unsafeToBytes(key))]
		runtime.KeepAlive(key)
		return e.value, ok
	}
	return
}

func (m *intMap[T]) contains(key *big.Int) bool {
	_, ok := m.get(key)
	return ok
}

func (m *intMap[T]) first() (ret T) {
	for _, v := range m.m {
		ret = v.value
		return
	}
	return
}

func (m *intMap[T]) forEach(f func(key *big.Int, value T) bool) {
	for _, v := range m.m {
		if !f(v.key, v.value) {
			return
		}
	}
}

func (m *intMap[T]) len() int {
	return len(m.m)
}
