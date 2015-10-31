package phoenix

// #include "go-clang.h"
import "C"

import (
	"fmt"
)

// 'Qualifiers' written next to the return and parameter types in ObjC method declarations.
type ObjCDeclQualifierKind uint32

const (
	ObjCDeclQualifier_None   ObjCDeclQualifierKind = C.CXObjCDeclQualifier_None
	ObjCDeclQualifier_In                           = C.CXObjCDeclQualifier_In
	ObjCDeclQualifier_Inout                        = C.CXObjCDeclQualifier_Inout
	ObjCDeclQualifier_Out                          = C.CXObjCDeclQualifier_Out
	ObjCDeclQualifier_Bycopy                       = C.CXObjCDeclQualifier_Bycopy
	ObjCDeclQualifier_Byref                        = C.CXObjCDeclQualifier_Byref
	ObjCDeclQualifier_Oneway                       = C.CXObjCDeclQualifier_Oneway
)

func (ocdqk ObjCDeclQualifierKind) Spelling() string {
	switch ocdqk {
	case ObjCDeclQualifier_None:
		return "ObjCDeclQualifier=None"
	case ObjCDeclQualifier_In:
		return "ObjCDeclQualifier=In"
	case ObjCDeclQualifier_Inout:
		return "ObjCDeclQualifier=Inout"
	case ObjCDeclQualifier_Out:
		return "ObjCDeclQualifier=Out"
	case ObjCDeclQualifier_Bycopy:
		return "ObjCDeclQualifier=Bycopy"
	case ObjCDeclQualifier_Byref:
		return "ObjCDeclQualifier=Byref"
	case ObjCDeclQualifier_Oneway:
		return "ObjCDeclQualifier=Oneway"

	}

	return fmt.Sprintf("ObjCDeclQualifierKind unkown %d", int(ocdqk))
}

func (ocdqk ObjCDeclQualifierKind) String() string {
	return ocdqk.Spelling()
}
