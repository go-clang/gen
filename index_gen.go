package phoenix

// #include "go-clang.h"
import "C"

// An "index" that consists of a set of translation units that would typically be linked together into an executable or library.
type Index struct {
	c C.CXIndex
}

// Destroy the given index. The index must not be destroyed until all of the translation units created within that index have been destroyed.
func (i Index) Dispose() {
	C.clang_disposeIndex(i.c)
}

// Gets the general options associated with a CXIndex. \returns A bitmask of options, a bitwise OR of CXGlobalOpt_XXX flags that are associated with the given CXIndex object.
func (i Index) GlobalOptions() uint16 {
	return uint16(C.clang_CXIndex_getGlobalOptions(i.c))
}

// An indexing action/session, to be applied to one or multiple translation units. \param CIdx The index object with which the index action will be associated.
func (i Index) Action_create() IndexAction {
	return IndexAction{C.clang_IndexAction_create(i.c)}
}
