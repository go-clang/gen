package phoenix

// #include "go-clang.h"
import "C"

type TUResourceUsageEntry struct {
	c C.CXTUResourceUsageEntry
}

func (turue TUResourceUsageEntry) Kind() TUResourceUsageKind {
	value := TUResourceUsageKind(turue.c.kind)
	return value
}

func (turue TUResourceUsageEntry) Amount() uint32 {
	value := uint32(turue.c.amount)
	return value
}
