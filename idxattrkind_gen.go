package phoenix

// #include "go-clang.h"
import "C"

type IdxAttrKind uint32

const (
	IdxAttr_Unexposed          IdxAttrKind = C.CXIdxAttr_Unexposed
	IdxAttr_IBAction                       = C.CXIdxAttr_IBAction
	IdxAttr_IBOutlet                       = C.CXIdxAttr_IBOutlet
	IdxAttr_IBOutletCollection             = C.CXIdxAttr_IBOutletCollection
)
