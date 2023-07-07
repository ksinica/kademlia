package kademlia

import (
	"math/big"
	"reflect"
	"runtime"
	"unsafe"
)

func unsafeSizeOf[T any]() int {
	var t T
	return int(unsafe.Sizeof(t))
}

func unsafeCast[T any](s []T) (ret []byte) {
	ret = make([]byte, 0)

	p := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&s)).Data)

	h := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
	h.Data = uintptr(p)
	h.Len = len(s) * unsafeSizeOf[T]()
	runtime.KeepAlive(s)
	return
}

func unsafeToBytes(z *big.Int) []byte {
	return unsafeCast(z.Bits())
}
