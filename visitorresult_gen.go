package phoenix

// #include "go-clang.h"
import "C"

// \defgroup CINDEX_HIGH Higher level API functions @{
type VisitorResult int

const (
	Visit_Break    VisitorResult = C.CXVisit_Break
	Visit_Continue               = C.CXVisit_Continue
)
