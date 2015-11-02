package phoenix

// #include "go-clang.h"
import "C"
import "fmt"

// Property attributes for a CXCursor_ObjCPropertyDecl.
type ObjCPropertyAttrKind uint32

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

func (ocpak ObjCPropertyAttrKind) Spelling() string {
	switch ocpak {
	case ObjCPropertyAttr_noattr:
		return "ObjCPropertyAttr=noattr"
	case ObjCPropertyAttr_readonly:
		return "ObjCPropertyAttr=readonly"
	case ObjCPropertyAttr_getter:
		return "ObjCPropertyAttr=getter"
	case ObjCPropertyAttr_assign:
		return "ObjCPropertyAttr=assign"
	case ObjCPropertyAttr_readwrite:
		return "ObjCPropertyAttr=readwrite"
	case ObjCPropertyAttr_retain:
		return "ObjCPropertyAttr=retain"
	case ObjCPropertyAttr_copy:
		return "ObjCPropertyAttr=copy"
	case ObjCPropertyAttr_nonatomic:
		return "ObjCPropertyAttr=nonatomic"
	case ObjCPropertyAttr_setter:
		return "ObjCPropertyAttr=setter"
	case ObjCPropertyAttr_atomic:
		return "ObjCPropertyAttr=atomic"
	case ObjCPropertyAttr_weak:
		return "ObjCPropertyAttr=weak"
	case ObjCPropertyAttr_strong:
		return "ObjCPropertyAttr=strong"
	case ObjCPropertyAttr_unsafe_unretained:
		return "ObjCPropertyAttr=unsafe_unretained"

	}

	return fmt.Sprintf("ObjCPropertyAttrKind unkown %d", int(ocpak))
}

func (ocpak ObjCPropertyAttrKind) String() string {
	return ocpak.Spelling()
}
