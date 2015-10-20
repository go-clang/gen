package phoenix

// #include "go-clang.h"
import "C"

// A comment AST node.
type Comment struct {
	c C.CXComment
}

// \param Comment a \c CXComment_Text AST node. \returns text contained in the AST node.
func (c Comment) TextComment_getText() string {
	cstr := cxstring{C.clang_TextComment_getText(c.c)}
	defer cstr.Dispose()

	return cstr.String()
}

// \param Comment a \c CXComment_InlineCommand AST node. \returns name of the inline command.
func (c Comment) InlineCommandComment_getCommandName() string {
	cstr := cxstring{C.clang_InlineCommandComment_getCommandName(c.c)}
	defer cstr.Dispose()

	return cstr.String()
}

// \param Comment a \c CXComment_HTMLStartTag or \c CXComment_HTMLEndTag AST node. \returns HTML tag name.
func (c Comment) HTMLTagComment_getTagName() string {
	cstr := cxstring{C.clang_HTMLTagComment_getTagName(c.c)}
	defer cstr.Dispose()

	return cstr.String()
}

// \param Comment a \c CXComment_BlockCommand AST node. \returns name of the block command.
func (c Comment) BlockCommandComment_getCommandName() string {
	cstr := cxstring{C.clang_BlockCommandComment_getCommandName(c.c)}
	defer cstr.Dispose()

	return cstr.String()
}

// \param Comment a \c CXComment_ParamCommand AST node. \returns parameter name.
func (c Comment) ParamCommandComment_getParamName() string {
	cstr := cxstring{C.clang_ParamCommandComment_getParamName(c.c)}
	defer cstr.Dispose()

	return cstr.String()
}

// \param Comment a \c CXComment_TParamCommand AST node. \returns template parameter name.
func (c Comment) TParamCommandComment_getParamName() string {
	cstr := cxstring{C.clang_TParamCommandComment_getParamName(c.c)}
	defer cstr.Dispose()

	return cstr.String()
}

// \param Comment a \c CXComment_VerbatimBlockLine AST node. \returns text contained in the AST node.
func (c Comment) VerbatimBlockLineComment_getText() string {
	cstr := cxstring{C.clang_VerbatimBlockLineComment_getText(c.c)}
	defer cstr.Dispose()

	return cstr.String()
}

// \param Comment a \c CXComment_VerbatimLine AST node. \returns text contained in the AST node.
func (c Comment) VerbatimLineComment_getText() string {
	cstr := cxstring{C.clang_VerbatimLineComment_getText(c.c)}
	defer cstr.Dispose()

	return cstr.String()
}

// Convert an HTML tag AST node to string. \param Comment a \c CXComment_HTMLStartTag or \c CXComment_HTMLEndTag AST node. \returns string containing an HTML tag.
func (c Comment) HTMLTagComment_getAsString() string {
	cstr := cxstring{C.clang_HTMLTagComment_getAsString(c.c)}
	defer cstr.Dispose()

	return cstr.String()
}

// Convert a given full parsed comment to an HTML fragment. Specific details of HTML layout are subject to change. Don't try to parse this HTML back into an AST, use other APIs instead. Currently the following CSS classes are used: \li "para-brief" for \ paragraph and equivalent commands; \li "para-returns" for \\returns paragraph and equivalent commands; \li "word-returns" for the "Returns" word in \\returns paragraph. Function argument documentation is rendered as a \<dl\> list with arguments sorted in function prototype order. CSS classes used: \li "param-name-index-NUMBER" for parameter name (\<dt\>); \li "param-descr-index-NUMBER" for parameter description (\<dd\>); \li "param-name-index-invalid" and "param-descr-index-invalid" are used if parameter index is invalid. Template parameter documentation is rendered as a \<dl\> list with parameters sorted in template parameter list order. CSS classes used: \li "tparam-name-index-NUMBER" for parameter name (\<dt\>); \li "tparam-descr-index-NUMBER" for parameter description (\<dd\>); \li "tparam-name-index-other" and "tparam-descr-index-other" are used for names inside template template parameters; \li "tparam-name-index-invalid" and "tparam-descr-index-invalid" are used if parameter position is invalid. \param Comment a \c CXComment_FullComment AST node. \returns string containing an HTML fragment.
func (c Comment) FullComment_getAsHTML() string {
	cstr := cxstring{C.clang_FullComment_getAsHTML(c.c)}
	defer cstr.Dispose()

	return cstr.String()
}

// Convert a given full parsed comment to an XML document. A Relax NG schema for the XML can be found in comment-xml-schema.rng file inside clang source tree. \param Comment a \c CXComment_FullComment AST node. \returns string containing an XML document.
func (c Comment) FullComment_getAsXML() string {
	cstr := cxstring{C.clang_FullComment_getAsXML(c.c)}
	defer cstr.Dispose()

	return cstr.String()
}
