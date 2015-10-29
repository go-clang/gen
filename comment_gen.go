package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// A comment AST node.
type Comment struct {
	c C.CXComment
}

func (c Comment) ASTNode() unsafe.Pointer {
	value := unsafe.Pointer(c.c.ASTNode)
	return value
}

func (c Comment) TranslationUnit() TranslationUnit {
	value := TranslationUnit{c.c.TranslationUnit}
	return value
}

// \param Comment AST node of any kind. \param ChildIdx child index (zero-based). \returns the specified child of the AST node.
func (c Comment) Child(ChildIdx uint16) Comment {
	return Comment{C.clang_Comment_getChild(c.c, C.uint(ChildIdx))}
}

// \param Comment a \c CXComment_InlineCommand AST node. \param ArgIdx argument index (zero-based). \returns text of the specified argument.
func (c Comment) InlineCommandComment_getArgText(ArgIdx uint16) string {
	o := cxstring{C.clang_InlineCommandComment_getArgText(c.c, C.uint(ArgIdx))}
	defer o.Dispose()

	return o.String()
}

// \param Comment a \c CXComment_HTMLStartTag AST node. \param AttrIdx attribute index (zero-based). \returns name of the specified attribute.
func (c Comment) HTMLStartTag_getAttrName(AttrIdx uint16) string {
	o := cxstring{C.clang_HTMLStartTag_getAttrName(c.c, C.uint(AttrIdx))}
	defer o.Dispose()

	return o.String()
}

// \param Comment a \c CXComment_HTMLStartTag AST node. \param AttrIdx attribute index (zero-based). \returns value of the specified attribute.
func (c Comment) HTMLStartTag_getAttrValue(AttrIdx uint16) string {
	o := cxstring{C.clang_HTMLStartTag_getAttrValue(c.c, C.uint(AttrIdx))}
	defer o.Dispose()

	return o.String()
}

// \param Comment a \c CXComment_BlockCommand AST node. \param ArgIdx argument index (zero-based). \returns text of the specified word-like argument.
func (c Comment) BlockCommandComment_getArgText(ArgIdx uint16) string {
	o := cxstring{C.clang_BlockCommandComment_getArgText(c.c, C.uint(ArgIdx))}
	defer o.Dispose()

	return o.String()
}

// \param Comment a \c CXComment_TParamCommand AST node. \returns zero-based parameter index in the template parameter list at a given nesting depth. For example, \verbatim template<typename C, template<typename T> class TT> void test(TT<int> aaa); \endverbatim for C and TT nesting depth is 0, so we can ask for index at depth 0: at depth 0 C's index is 0, TT's index is 1. For T nesting depth is 1, so we can ask for index at depth 0 and 1: at depth 0 T's index is 1 (same as TT's), at depth 1 T's index is 0.
func (c Comment) TParamCommandComment_getIndex(Depth uint16) uint16 {
	return uint16(C.clang_TParamCommandComment_getIndex(c.c, C.uint(Depth)))
}
