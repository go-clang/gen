package phoenix

// #include "go-clang.h"
import "C"

// Flags that control the reparsing of translation units. The enumerators in this enumeration type are meant to be bitwise ORed together to specify which options should be used when reparsing the translation unit.
type Reparse_Flags int

const (
	// Used to indicate that no special reparsing options are needed.
	Reparse_None Reparse_Flags = C.CXReparse_None
)
