package phoenix

// #include "./clang-c/CXString.h"
// #include "go-clang.h"
import "C"

// A character string. The \c CXString type is used to return strings from the interface when the ownership of that string might different from one call to the next. Use \c clang_getCString() to retrieve the string data and, once finished with the string data, call \c clang_disposeString() to free the string.
type String struct {
	c C.CXString
}

func (s String) Private_flags() uint16 {
	value := uint16(s.c.private_flags)
	return value
}
