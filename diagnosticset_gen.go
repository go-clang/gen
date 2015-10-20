package phoenix

// #include "go-clang.h"
import "C"

// A group of CXDiagnostics.
type DiagnosticSet struct {
	c C.CXDiagnosticSet
}
