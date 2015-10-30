package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// The memory usage of a CXTranslationUnit, broken into categories.
type TUResourceUsage struct {
	c C.CXTUResourceUsage
}

func (turu TUResourceUsage) Data() unsafe.Pointer {
	value := unsafe.Pointer(turu.c.data)
	return value
}

func (turu TUResourceUsage) NumEntries() uint16 {
	value := uint16(turu.c.numEntries)
	return value
}

func (turu TUResourceUsage) Entries() *TUResourceUsageEntry {
	value := TUResourceUsageEntry{*turu.c.entries}
	return &value
}

func (turu TUResourceUsage) Dispose() {
	C.clang_disposeCXTUResourceUsage(turu.c)
}
