package phoenix

// #include "go-clang.h"
import "C"

// 'Qualifiers' written next to the return and parameter types in ObjC method declarations.
type ObjCDeclQualifierKind int

const (
	ObjCDeclQualifier_None   ObjCDeclQualifierKind = C.CXObjCDeclQualifier_None
	ObjCDeclQualifier_In                           = C.CXObjCDeclQualifier_In
	ObjCDeclQualifier_Inout                        = C.CXObjCDeclQualifier_Inout
	ObjCDeclQualifier_Out                          = C.CXObjCDeclQualifier_Out
	ObjCDeclQualifier_Bycopy                       = C.CXObjCDeclQualifier_Bycopy
	ObjCDeclQualifier_Byref                        = C.CXObjCDeclQualifier_Byref
	ObjCDeclQualifier_Oneway                       = C.CXObjCDeclQualifier_Oneway
)
