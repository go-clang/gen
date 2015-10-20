package phoenix

// #include "go-clang.h"
import "C"

// Extra C++ template information for an entity. This can apply to: CXIdxEntity_Function CXIdxEntity_CXXClass CXIdxEntity_CXXStaticMethod CXIdxEntity_CXXInstanceMethod CXIdxEntity_CXXConstructor CXIdxEntity_CXXConversionFunction CXIdxEntity_CXXTypeAlias
type IdxEntityCXXTemplateKind int

const (
	IdxEntity_NonTemplate                   IdxEntityCXXTemplateKind = C.CXIdxEntity_NonTemplate
	IdxEntity_Template                                               = C.CXIdxEntity_Template
	IdxEntity_TemplatePartialSpecialization                          = C.CXIdxEntity_TemplatePartialSpecialization
	IdxEntity_TemplateSpecialization                                 = C.CXIdxEntity_TemplateSpecialization
)
