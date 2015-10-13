package phoenix

// #include "go-clang.h"
import "C"

//
type RefQualifierKind int

const (
	// No ref-qualifier was provided.
	RefQualifier_None RefQualifierKind = C.CXRefQualifier_None
	// An lvalue ref-qualifier was provided (\c &).
	RefQualifier_LValue = C.CXRefQualifier_LValue
	// An rvalue ref-qualifier was provided (\c &&).
	RefQualifier_RValue = C.CXRefQualifier_RValue
)
