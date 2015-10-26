package clang

// #include <stdlib.h>
// #include "go-clang.h"
//
import "C"

/**
 * \brief A comment AST node.
 */
type Comment struct {
	c C.CXComment
}

/**
 * \param Comment AST node of any kind.
 *
 * \param ChildIdx child index (zero-based).
 *
 * \returns the specified child of the AST node.
 */
func (c Comment) Child(idx int) Comment {
	return Comment{C.clang_Comment_getChild(c.c, C.unsigned(idx))}
}
