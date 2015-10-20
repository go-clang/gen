package phoenix

// #include "go-clang.h"
import "C"

// Property attributes for a \c CXCursor_ObjCPropertyDecl.
type ObjCPropertyAttrKind int

const (
	ObjCPropertyAttr_noattr            ObjCPropertyAttrKind = C.CXObjCPropertyAttr_noattr
	ObjCPropertyAttr_readonly                               = C.CXObjCPropertyAttr_readonly
	ObjCPropertyAttr_getter                                 = C.CXObjCPropertyAttr_getter
	ObjCPropertyAttr_assign                                 = C.CXObjCPropertyAttr_assign
	ObjCPropertyAttr_readwrite                              = C.CXObjCPropertyAttr_readwrite
	ObjCPropertyAttr_retain                                 = C.CXObjCPropertyAttr_retain
	ObjCPropertyAttr_copy                                   = C.CXObjCPropertyAttr_copy
	ObjCPropertyAttr_nonatomic                              = C.CXObjCPropertyAttr_nonatomic
	ObjCPropertyAttr_setter                                 = C.CXObjCPropertyAttr_setter
	ObjCPropertyAttr_atomic                                 = C.CXObjCPropertyAttr_atomic
	ObjCPropertyAttr_weak                                   = C.CXObjCPropertyAttr_weak
	ObjCPropertyAttr_strong                                 = C.CXObjCPropertyAttr_strong
	ObjCPropertyAttr_unsafe_unretained                      = C.CXObjCPropertyAttr_unsafe_unretained
)
