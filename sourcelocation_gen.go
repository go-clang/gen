package phoenix

// #include "go-clang.h"
import "C"

// Identifies a specific source location within a translation unit. Use clang_getExpansionLocation() or clang_getSpellingLocation() to map a source location to a particular file, line, and column.
type SourceLocation struct {
	c C.CXSourceLocation
}
