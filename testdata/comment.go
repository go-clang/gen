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
 * \returns number of children of the AST node.
 */
func (c Comment) NumChildren() int {
	return int(C.clang_Comment_getNumChildren(c.c))
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

/**
 * \returns non-zero if \c Comment is inline content and has a newline
 * immediately following it in the comment text.  Newlines between paragraphs
 * do not count.
 */
func (c Comment) HasTrailingNewline() bool {
	o := C.clang_InlineContentComment_hasTrailingNewline(c.c)
	if 0 != o {
		return true
	}
	return false
}
