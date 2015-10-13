package phoenix

// #include "go-clang.h"
import "C"

type RefQualifierKind int

const (
	/** \brief No ref-qualifier was provided. */
	RefQualifier_None RefQualifierKind = C.CXRefQualifier_None
	/** \brief An lvalue ref-qualifier was provided (\c &). */
	RefQualifier_LValue RefQualifierKind = C.CXRefQualifier_LValue
	/** \brief An rvalue ref-qualifier was provided (\c &&). */
	RefQualifier_RValue RefQualifierKind = C.CXRefQualifier_RValue
)
