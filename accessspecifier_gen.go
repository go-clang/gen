package phoenix

// #include "go-clang.h"
import "C"

// Represents the C++ access control level to a base class for a cursor with kind CX_CXXBaseSpecifier.
type AccessSpecifier int

const (
	InvalidAccessSpecifier AccessSpecifier = C.CX_CXXInvalidAccessSpecifier
	Public                                 = C.CX_CXXPublic
	Protected                              = C.CX_CXXProtected
	Private                                = C.CX_CXXPrivate
)
