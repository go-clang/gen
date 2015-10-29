package phoenix

// #include "go-clang.h"
import "C"

// An indexing action/session, to be applied to one or multiple translation units.
type IndexAction struct {
	c C.CXIndexAction
}
