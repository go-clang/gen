package phoenix

// #include "go-clang.h"
import "C"

// A comment AST node.
type Comment struct {
	c C.CXComment
}
