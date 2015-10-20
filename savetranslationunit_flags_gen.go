package phoenix

// #include "go-clang.h"
import "C"

// Flags that control how translation units are saved. The enumerators in this enumeration type are meant to be bitwise ORed together to specify which options should be used when saving the translation unit.
type SaveTranslationUnit_Flags int

const (
	// Used to indicate that no special saving options are needed.
	SaveTranslationUnit_None SaveTranslationUnit_Flags = C.CXSaveTranslationUnit_None
)
