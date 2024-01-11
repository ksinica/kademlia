package kademlia

import (
	"math/big"
	"runtime"
	"unsafe"
)

func unsafeSizeOf[T any]() int {
	var t T
	return int(unsafe.Sizeof(t))
}

func unsafeCast[T any](s []T) (ret string) {
	p := unsafe.SliceData(s)
	ret = unsafe.String(
		(*byte)(unsafe.Pointer(p)),
		unsafeSizeOf[T](),
	)
	runtime.KeepAlive(s)
	return
}

func unsafeBigIntToString(z *big.Int) string {
	return unsafeCast(z.Bits())
}
