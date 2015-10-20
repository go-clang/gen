package phoenix

// #include "go-clang.h"
import "C"

type IdxEntityKind int

const (
	IdxEntity_Unexposed             IdxEntityKind = C.CXIdxEntity_Unexposed
	IdxEntity_Typedef                             = C.CXIdxEntity_Typedef
	IdxEntity_Function                            = C.CXIdxEntity_Function
	IdxEntity_Variable                            = C.CXIdxEntity_Variable
	IdxEntity_Field                               = C.CXIdxEntity_Field
	IdxEntity_EnumConstant                        = C.CXIdxEntity_EnumConstant
	IdxEntity_ObjCClass                           = C.CXIdxEntity_ObjCClass
	IdxEntity_ObjCProtocol                        = C.CXIdxEntity_ObjCProtocol
	IdxEntity_ObjCCategory                        = C.CXIdxEntity_ObjCCategory
	IdxEntity_ObjCInstanceMethod                  = C.CXIdxEntity_ObjCInstanceMethod
	IdxEntity_ObjCClassMethod                     = C.CXIdxEntity_ObjCClassMethod
	IdxEntity_ObjCProperty                        = C.CXIdxEntity_ObjCProperty
	IdxEntity_ObjCIvar                            = C.CXIdxEntity_ObjCIvar
	IdxEntity_Enum                                = C.CXIdxEntity_Enum
	IdxEntity_Struct                              = C.CXIdxEntity_Struct
	IdxEntity_Union                               = C.CXIdxEntity_Union
	IdxEntity_CXXClass                            = C.CXIdxEntity_CXXClass
	IdxEntity_CXXNamespace                        = C.CXIdxEntity_CXXNamespace
	IdxEntity_CXXNamespaceAlias                   = C.CXIdxEntity_CXXNamespaceAlias
	IdxEntity_CXXStaticVariable                   = C.CXIdxEntity_CXXStaticVariable
	IdxEntity_CXXStaticMethod                     = C.CXIdxEntity_CXXStaticMethod
	IdxEntity_CXXInstanceMethod                   = C.CXIdxEntity_CXXInstanceMethod
	IdxEntity_CXXConstructor                      = C.CXIdxEntity_CXXConstructor
	IdxEntity_CXXDestructor                       = C.CXIdxEntity_CXXDestructor
	IdxEntity_CXXConversionFunction               = C.CXIdxEntity_CXXConversionFunction
	IdxEntity_CXXTypeAlias                        = C.CXIdxEntity_CXXTypeAlias
	IdxEntity_CXXInterface                        = C.CXIdxEntity_CXXInterface
)
