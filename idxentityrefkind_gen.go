package phoenix

// #include "go-clang.h"
import "C"

// Data for IndexerCallbacks#indexEntityReference.
type IdxEntityRefKind int

const (
	// The entity is referenced directly in user's code.
	IdxEntityRef_Direct IdxEntityRefKind = C.CXIdxEntityRef_Direct
	// An implicit reference, e.g. a reference of an ObjC method via the dot syntax.
	IdxEntityRef_Implicit = C.CXIdxEntityRef_Implicit
)
