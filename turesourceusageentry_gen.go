package phoenix

// #include "go-clang.h"
import "C"

type TUResourceUsageEntry struct {
	c C.CXTUResourceUsageEntry
}

func (turue TUResourceUsageEntry) Kind() TUResourceUsageKind {
	return TUResourceUsageKind(turue.c.kind)
}

func (turue TUResourceUsageEntry) Amount() uint32 {
	return uint32(turue.c.amount)
}
