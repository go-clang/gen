package phoenix

// #include "go-clang.h"
import "C"

// A single result of code completion.
type CompletionResult struct {
	c C.CXCompletionResult
}
