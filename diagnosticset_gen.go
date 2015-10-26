package phoenix

// #include "go-clang.h"
import "C"

// A group of CXDiagnostics.
type DiagnosticSet struct {
	c C.CXDiagnosticSet
}

// Determine the number of diagnostics in a CXDiagnosticSet.
func (ds DiagnosticSet) NumDiagnosticsInSet() uint16 {
	return uint16(C.clang_getNumDiagnosticsInSet(ds.c))
}

// Release a CXDiagnosticSet and all of its contained diagnostics.
func (ds DiagnosticSet) Dispose() {
	C.clang_disposeDiagnosticSet(ds.c)
}
