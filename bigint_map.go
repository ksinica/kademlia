package kademlia

import (
	"math/big"
	"runtime"
)

type bigIntMapEntry[T any] struct {
	key   *big.Int
	value T
}

// bigIntMap is a map with big.Int keys. It offers a convenient and optimized way
// to relate some arbitrary types with a big.Int.
type bigIntMap[T any] struct {
	m map[string]bigIntMapEntry[T]
}

// set sets the value for a key.
func (m *bigIntMap[T]) set(key *big.Int, value T) {
	if m.m == nil {
		m.m = make(map[string]bigIntMapEntry[T])
	}
	m.m[string(unsafeToBytes(key))] = bigIntMapEntry[T]{
		key:   key,
		value: value,
	}
	runtime.KeepAlive(key)
}

// get returns the value stored in the map for a key, or zero value if no value
// is present. The ok result indicates whether value was found in the map.
func (m *bigIntMap[T]) get(key *big.Int) (value T, ok bool) {
	if m.m != nil {
		e, ok := m.m[string(unsafeToBytes(key))]
		runtime.KeepAlive(key)
		return e.value, ok
	}
	return
}

// contains returns true if the map contains some value associated with
// a given key.
func (m *bigIntMap[T]) contains(key *big.Int) bool {
	_, ok := m.get(key)
	return ok
}

func (m *bigIntMap[T]) first() (ret T) {
	for _, v := range m.m {
		ret = v.value
		return
	}
	return
}

// forEach calls f sequentially for each key and value present in the map.
// If f returns false, it stops the iteration.
func (m *bigIntMap[T]) forEach(f func(key *big.Int, value T) bool) {
	for _, v := range m.m {
		if !f(v.key, v.value) {
			return
		}
	}
}

func (m *bigIntMap[T]) len() int {
	return len(m.m)
}
