package kademlia

import (
	"math/big"
	"reflect"
	"runtime"
	"unsafe"
)

func bigWordToBytes(s []big.Word) (ret []byte) {
	const (
		W = int(unsafe.Sizeof(big.Word(0)))
	)
	ret = make([]byte, 0)

	p := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&s)).Data)

	h := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
	h.Data = uintptr(p)
	h.Len = len(s) * W
	runtime.KeepAlive(s)
	return
}

func bigIntToBytes(z *big.Int) []byte {
	return bigWordToBytes(z.Bits())
}

// intMap is a map with big.Int keys. It offers a convenient and optimized way
// to relate some arbitrary types with a big.Int.
type intMapEntry[T any] struct {
	key   *big.Int
	value T
}

type intMap[T any] struct {
	m map[string]intMapEntry[T]
}

func (m *intMap[T]) set(key *big.Int, value T) {
	if m.m == nil {
		m.m = make(map[string]intMapEntry[T])
	}
	m.m[string(bigIntToBytes(key))] = intMapEntry[T]{key: key, value: value}
}

func (m *intMap[T]) get(key *big.Int) (value T, ok bool) {
	if m.m != nil {
		e, ok := m.m[string(bigIntToBytes(key))]
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

func (m *intMap[T]) length() int {
	return len(m.m)
}
