package phoenix

// #include "go-clang.h"
import "C"

// Describe the "language" of the entity referred to by a cursor.
type LanguageKind uint32

const (
	Language_Invalid   LanguageKind = C.CXLanguage_Invalid
	Language_C                      = C.CXLanguage_C
	Language_ObjC                   = C.CXLanguage_ObjC
	Language_CPlusPlus              = C.CXLanguage_CPlusPlus
)
