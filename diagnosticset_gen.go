package phoenix

// #include "go-clang.h"
import "C"

// A group of CXDiagnostics.
type DiagnosticSet struct {
	c C.CXDiagnosticSet
}

// Release a CXDiagnosticSet and all of its contained diagnostics.
func (ds DiagnosticSet) Dispose() {
	C.clang_disposeDiagnosticSet(ds.c)
}
