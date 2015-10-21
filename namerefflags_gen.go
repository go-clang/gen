package phoenix

// #include "go-clang.h"
import "C"

type NameRefFlags uint32

const (
	// Include the nested-name-specifier, e.g. Foo:: in x.Foo::y, in the range.
	NameRange_WantQualifier NameRefFlags = C.CXNameRange_WantQualifier
	// Include the explicit template arguments, e.g. \<int> in x.f<int>, in the range.
	NameRange_WantTemplateArgs = C.CXNameRange_WantTemplateArgs
	/*
	 * \brief If the name is non-contiguous, return the full spanning range.
	 *
	 * Non-contiguous names occur in Objective-C when a selector with two or more
	 * parameters is used, or in C++ when using an operator:
	 * \code
	 * [object doSomething:here withValue:there]; // ObjC
	 * return some_vector[1]; // C++
	 * \endcode
	 */
	NameRange_WantSinglePiece = C.CXNameRange_WantSinglePiece
)
