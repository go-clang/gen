package phoenix

// #include "go-clang.h"
import "C"

// Flags that can be passed to \c clang_codeCompleteAt() to modify its behavior. The enumerators in this enumeration can be bitwise-OR'd together to provide multiple options to \c clang_codeCompleteAt().
type CodeComplete_Flags int

const (
	// Whether to include macros within the set of code completions returned.
	CodeComplete_IncludeMacros CodeComplete_Flags = C.CXCodeComplete_IncludeMacros
	// Whether to include code patterns for language constructs within the set of code completions, e.g., for loops.
	CodeComplete_IncludeCodePatterns = C.CXCodeComplete_IncludeCodePatterns
	// Whether to include brief documentation within the set of code completions returned.
	CodeComplete_IncludeBriefComments = C.CXCodeComplete_IncludeBriefComments
)
