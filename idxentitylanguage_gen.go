package phoenix

// #include "go-clang.h"
import "C"

type IdxEntityLanguage int

const (
	IdxEntityLang_None IdxEntityLanguage = C.CXIdxEntityLang_None
	IdxEntityLang_C                      = C.CXIdxEntityLang_C
	IdxEntityLang_ObjC                   = C.CXIdxEntityLang_ObjC
	IdxEntityLang_CXX                    = C.CXIdxEntityLang_CXX
)
