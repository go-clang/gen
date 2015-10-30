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

// \param Comment AST node of any kind. \returns the type of the AST node.
func (c Comment) Kind() CommentKind {
	return CommentKind(C.clang_Comment_getKind(c.c))
}

// \param Comment AST node of any kind. \returns number of children of the AST node.
func (c Comment) NumChildren() uint16 {
	return uint16(C.clang_Comment_getNumChildren(c.c))
}

// \param Comment AST node of any kind. \param ChildIdx child index (zero-based). \returns the specified child of the AST node.
func (c Comment) Child(ChildIdx uint16) Comment {
	return Comment{C.clang_Comment_getChild(c.c, C.uint(ChildIdx))}
}

// A \c CXComment_Paragraph node is considered whitespace if it contains only \c CXComment_Text nodes that are empty or whitespace. Other AST nodes (except \c CXComment_Paragraph and \c CXComment_Text) are never considered whitespace. \returns non-zero if \c Comment is whitespace.
func (c Comment) IsWhitespace() bool {
	o := C.clang_Comment_isWhitespace(c.c)

	return o != C.uint(0)
}

// \returns non-zero if \c Comment is inline content and has a newline immediately following it in the comment text. Newlines between paragraphs do not count.
func (c Comment) InlineContentComment_HasTrailingNewline() bool {
	o := C.clang_InlineContentComment_hasTrailingNewline(c.c)

	return o != C.uint(0)
}

// \param Comment a \c CXComment_Text AST node. \returns text contained in the AST node.
func (c Comment) TextComment_getText() string {
	o := cxstring{C.clang_TextComment_getText(c.c)}
	defer o.Dispose()

	return o.String()
}

// \param Comment a \c CXComment_InlineCommand AST node. \returns name of the inline command.
func (c Comment) InlineCommandComment_getCommandName() string {
	o := cxstring{C.clang_InlineCommandComment_getCommandName(c.c)}
	defer o.Dispose()

	return o.String()
}

// \param Comment a \c CXComment_InlineCommand AST node. \returns the most appropriate rendering mode, chosen on command semantics in Doxygen.
func (c Comment) InlineCommandComment_getRenderKind() CommentInlineCommandRenderKind {
	return CommentInlineCommandRenderKind(C.clang_InlineCommandComment_getRenderKind(c.c))
}

// \param Comment a \c CXComment_InlineCommand AST node. \returns number of command arguments.
func (c Comment) InlineCommandComment_getNumArgs() uint16 {
	return uint16(C.clang_InlineCommandComment_getNumArgs(c.c))
}

// \param Comment a \c CXComment_InlineCommand AST node. \param ArgIdx argument index (zero-based). \returns text of the specified argument.
func (c Comment) InlineCommandComment_getArgText(ArgIdx uint16) string {
	o := cxstring{C.clang_InlineCommandComment_getArgText(c.c, C.uint(ArgIdx))}
	defer o.Dispose()

	return o.String()
}

// \param Comment a \c CXComment_HTMLStartTag or \c CXComment_HTMLEndTag AST node. \returns HTML tag name.
func (c Comment) HTMLTagComment_getTagName() string {
	o := cxstring{C.clang_HTMLTagComment_getTagName(c.c)}
	defer o.Dispose()

	return o.String()
}

// \param Comment a \c CXComment_HTMLStartTag AST node. \returns non-zero if tag is self-closing (for example, &lt;br /&gt;).
func (c Comment) HTMLStartTagComment_IsSelfClosing() bool {
	o := C.clang_HTMLStartTagComment_isSelfClosing(c.c)

	return o != C.uint(0)
}

// \param Comment a \c CXComment_HTMLStartTag AST node. \returns number of attributes (name-value pairs) attached to the start tag.
func (c Comment) HTMLStartTag_getNumAttrs() uint16 {
	return uint16(C.clang_HTMLStartTag_getNumAttrs(c.c))
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

// \param Comment a \c CXComment_BlockCommand AST node. \returns name of the block command.
func (c Comment) BlockCommandComment_getCommandName() string {
	o := cxstring{C.clang_BlockCommandComment_getCommandName(c.c)}
	defer o.Dispose()

	return o.String()
}

// \param Comment a \c CXComment_BlockCommand AST node. \returns number of word-like arguments.
func (c Comment) BlockCommandComment_getNumArgs() uint16 {
	return uint16(C.clang_BlockCommandComment_getNumArgs(c.c))
}

// \param Comment a \c CXComment_BlockCommand AST node. \param ArgIdx argument index (zero-based). \returns text of the specified word-like argument.
func (c Comment) BlockCommandComment_getArgText(ArgIdx uint16) string {
	o := cxstring{C.clang_BlockCommandComment_getArgText(c.c, C.uint(ArgIdx))}
	defer o.Dispose()

	return o.String()
}

// \param Comment a \c CXComment_BlockCommand or \c CXComment_VerbatimBlockCommand AST node. \returns paragraph argument of the block command.
func (c Comment) BlockCommandComment_getParagraph() Comment {
	return Comment{C.clang_BlockCommandComment_getParagraph(c.c)}
}

// \param Comment a \c CXComment_ParamCommand AST node. \returns parameter name.
func (c Comment) ParamCommandComment_getParamName() string {
	o := cxstring{C.clang_ParamCommandComment_getParamName(c.c)}
	defer o.Dispose()

	return o.String()
}

// \param Comment a \c CXComment_ParamCommand AST node. \returns non-zero if the parameter that this AST node represents was found in the function prototype and \c clang_ParamCommandComment_getParamIndex function will return a meaningful value.
func (c Comment) ParamCommandComment_IsParamIndexValid() bool {
	o := C.clang_ParamCommandComment_isParamIndexValid(c.c)

	return o != C.uint(0)
}

// \param Comment a \c CXComment_ParamCommand AST node. \returns zero-based parameter index in function prototype.
func (c Comment) ParamCommandComment_getParamIndex() uint16 {
	return uint16(C.clang_ParamCommandComment_getParamIndex(c.c))
}

// \param Comment a \c CXComment_ParamCommand AST node. \returns non-zero if parameter passing direction was specified explicitly in the comment.
func (c Comment) ParamCommandComment_IsDirectionExplicit() bool {
	o := C.clang_ParamCommandComment_isDirectionExplicit(c.c)

	return o != C.uint(0)
}

// \param Comment a \c CXComment_ParamCommand AST node. \returns parameter passing direction.
func (c Comment) ParamCommandComment_getDirection() CommentParamPassDirection {
	return CommentParamPassDirection(C.clang_ParamCommandComment_getDirection(c.c))
}

// \param Comment a \c CXComment_TParamCommand AST node. \returns template parameter name.
func (c Comment) TParamCommandComment_getParamName() string {
	o := cxstring{C.clang_TParamCommandComment_getParamName(c.c)}
	defer o.Dispose()

	return o.String()
}

// \param Comment a \c CXComment_TParamCommand AST node. \returns non-zero if the parameter that this AST node represents was found in the template parameter list and \c clang_TParamCommandComment_getDepth and \c clang_TParamCommandComment_getIndex functions will return a meaningful value.
func (c Comment) TParamCommandComment_IsParamPositionValid() bool {
	o := C.clang_TParamCommandComment_isParamPositionValid(c.c)

	return o != C.uint(0)
}

// \param Comment a \c CXComment_TParamCommand AST node. \returns zero-based nesting depth of this parameter in the template parameter list. For example, \verbatim template<typename C, template<typename T> class TT> void test(TT<int> aaa); \endverbatim for C and TT nesting depth is 0, for T nesting depth is 1.
func (c Comment) TParamCommandComment_getDepth() uint16 {
	return uint16(C.clang_TParamCommandComment_getDepth(c.c))
}

// \param Comment a \c CXComment_TParamCommand AST node. \returns zero-based parameter index in the template parameter list at a given nesting depth. For example, \verbatim template<typename C, template<typename T> class TT> void test(TT<int> aaa); \endverbatim for C and TT nesting depth is 0, so we can ask for index at depth 0: at depth 0 C's index is 0, TT's index is 1. For T nesting depth is 1, so we can ask for index at depth 0 and 1: at depth 0 T's index is 1 (same as TT's), at depth 1 T's index is 0.
func (c Comment) TParamCommandComment_getIndex(Depth uint16) uint16 {
	return uint16(C.clang_TParamCommandComment_getIndex(c.c, C.uint(Depth)))
}

// \param Comment a \c CXComment_VerbatimBlockLine AST node. \returns text contained in the AST node.
func (c Comment) VerbatimBlockLineComment_getText() string {
	o := cxstring{C.clang_VerbatimBlockLineComment_getText(c.c)}
	defer o.Dispose()

	return o.String()
}

// \param Comment a \c CXComment_VerbatimLine AST node. \returns text contained in the AST node.
func (c Comment) VerbatimLineComment_getText() string {
	o := cxstring{C.clang_VerbatimLineComment_getText(c.c)}
	defer o.Dispose()

	return o.String()
}

// Convert an HTML tag AST node to string. \param Comment a \c CXComment_HTMLStartTag or \c CXComment_HTMLEndTag AST node. \returns string containing an HTML tag.
func (c Comment) HTMLTagComment_getAsString() string {
	o := cxstring{C.clang_HTMLTagComment_getAsString(c.c)}
	defer o.Dispose()

	return o.String()
}

// Convert a given full parsed comment to an HTML fragment. Specific details of HTML layout are subject to change. Don't try to parse this HTML back into an AST, use other APIs instead. Currently the following CSS classes are used: \li "para-brief" for \ paragraph and equivalent commands; \li "para-returns" for \\returns paragraph and equivalent commands; \li "word-returns" for the "Returns" word in \\returns paragraph. Function argument documentation is rendered as a \<dl\> list with arguments sorted in function prototype order. CSS classes used: \li "param-name-index-NUMBER" for parameter name (\<dt\>); \li "param-descr-index-NUMBER" for parameter description (\<dd\>); \li "param-name-index-invalid" and "param-descr-index-invalid" are used if parameter index is invalid. Template parameter documentation is rendered as a \<dl\> list with parameters sorted in template parameter list order. CSS classes used: \li "tparam-name-index-NUMBER" for parameter name (\<dt\>); \li "tparam-descr-index-NUMBER" for parameter description (\<dd\>); \li "tparam-name-index-other" and "tparam-descr-index-other" are used for names inside template template parameters; \li "tparam-name-index-invalid" and "tparam-descr-index-invalid" are used if parameter position is invalid. \param Comment a \c CXComment_FullComment AST node. \returns string containing an HTML fragment.
func (c Comment) FullComment_getAsHTML() string {
	o := cxstring{C.clang_FullComment_getAsHTML(c.c)}
	defer o.Dispose()

	return o.String()
}

// Convert a given full parsed comment to an XML document. A Relax NG schema for the XML can be found in comment-xml-schema.rng file inside clang source tree. \param Comment a \c CXComment_FullComment AST node. \returns string containing an XML document.
func (c Comment) FullComment_getAsXML() string {
	o := cxstring{C.clang_FullComment_getAsXML(c.c)}
	defer o.Dispose()

	return o.String()
}
