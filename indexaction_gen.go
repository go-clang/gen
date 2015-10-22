package phoenix

// #include "go-clang.h"
import "C"

// An indexing action/session, to be applied to one or multiple translation units.
type IndexAction struct {
	c C.CXIndexAction
}

// Destroy the given index action. The index action must not be destroyed until all of the translation units created within that index action have been destroyed.
func (ia IndexAction) Dispose() {
	C.clang_IndexAction_dispose(ia.c)
}
