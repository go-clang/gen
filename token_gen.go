package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// Describes a single preprocessing token.
type Token struct {
	c C.CXToken
}

func (t Token) Int_data() []uint16 {
	sc := []uint16{}

	length := 4
	goslice := (*[1 << 30]C.uint)(unsafe.Pointer(&t.c.int_data))[:length:length]

	for is := 0; is < length; is++ {
		sc = append(sc, uint16(goslice[is]))
	}

	return sc
}

func (t Token) Ptr_data() unsafe.Pointer {
	value := unsafe.Pointer(t.c.ptr_data)
	return value
}
