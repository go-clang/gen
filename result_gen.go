package phoenix

// #include "go-clang.h"
import "C"

type Result uint32

const (
	// Function returned successfully.
	Result_Success Result = C.CXResult_Success
	// One of the parameters was invalid for the function.
	Result_Invalid = C.CXResult_Invalid
	// The function was terminated by a callback (e.g. it returned CXVisit_Break)
	Result_VisitBreak = C.CXResult_VisitBreak
)
