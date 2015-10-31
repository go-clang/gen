package phoenix

// #include "go-clang.h"
import "C"

import (
	"fmt"
)

// Represents the C++ access control level to a base class for a cursor with kind CX_CXXBaseSpecifier.
type AccessSpecifier uint32

const (
	InvalidAccessSpecifier AccessSpecifier = C.CX_CXXInvalidAccessSpecifier
	Public                                 = C.CX_CXXPublic
	Protected                              = C.CX_CXXProtected
	Private                                = C.CX_CXXPrivate
)

func (as AccessSpecifier) Spelling() string {
	switch as {
	case InvalidAccessSpecifier:
		return "InvalidAccessSpecifier"
	case Public:
		return "Public"
	case Protected:
		return "Protected"
	case Private:
		return "Private"

	}

	return fmt.Sprintf("AccessSpecifier unkown %d", int(as))
}

func (as AccessSpecifier) String() string {
	return as.Spelling()
}
