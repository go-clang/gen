package phoenix

// #include "go-clang.h"
import "C"

// The memory usage of a CXTranslationUnit, broken into categories.
type TUResourceUsage struct {
	c C.CXTUResourceUsage
}

func (turu TUResourceUsage) Dispose() {
	C.clang_disposeCXTUResourceUsage(turu.c)
}

func (turu TUResourceUsage) NumEntries() uint16 {
	return uint16(turu.c.numEntries)
}

func (turu TUResourceUsage) Entries() *TUResourceUsageEntry {
	o := *turu.c.entries

	return &TUResourceUsageEntry{o}
}
