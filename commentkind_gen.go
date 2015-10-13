package phoenix

// #include "go-clang.h"
import "C"

// Describes the type of the comment AST node (\c CXComment). A comment node can be considered block content (e. g., paragraph), inline content (plain text) or neither (the root AST node).
type CommentKind int

const (
	// Null comment. No AST node is constructed at the requested location because there is no text or a syntax error.
	Comment_Null CommentKind = C.CXComment_Null
	// Plain text. Inline content.
	Comment_Text = C.CXComment_Text
	// A command with word-like arguments that is considered inline content. For example: \\c command.
	Comment_InlineCommand = C.CXComment_InlineCommand
	// HTML start tag with attributes (name-value pairs). Considered inline content. For example: \verbatim <br> <br /> <a href="http://example.org/"> \endverbatim
	Comment_HTMLStartTag = C.CXComment_HTMLStartTag
	// HTML end tag. Considered inline content. For example: \verbatim </a> \endverbatim
	Comment_HTMLEndTag = C.CXComment_HTMLEndTag
	// A paragraph, contains inline comment. The paragraph itself is block content.
	Comment_Paragraph = C.CXComment_Paragraph
	// A command that has zero or more word-like arguments (number of word-like arguments depends on command name) and a paragraph as an argument. Block command is block content. Paragraph argument is also a child of the block command. For example: \ has 0 word-like arguments and a paragraph argument. AST nodes of special kinds that parser knows about (e. g., \\param command) have their own node kinds.
	Comment_BlockCommand = C.CXComment_BlockCommand
	// A \\param or \\arg command that describes the function parameter (name, passing direction, description). For example: \\param [in] ParamName description.
	Comment_ParamCommand = C.CXComment_ParamCommand
	// A \\tparam command that describes a template parameter (name and description). For example: \\tparam T description.
	Comment_TParamCommand = C.CXComment_TParamCommand
	// A verbatim block command (e. g., preformatted code). Verbatim block has an opening and a closing command and contains multiple lines of text (\c CXComment_VerbatimBlockLine child nodes). For example: \\verbatim aaa \\endverbatim
	Comment_VerbatimBlockCommand = C.CXComment_VerbatimBlockCommand
	// A line of text that is contained within a CXComment_VerbatimBlockCommand node.
	Comment_VerbatimBlockLine = C.CXComment_VerbatimBlockLine
	// A verbatim line command. Verbatim line has an opening command, a single line of text (up to the newline after the opening command) and has no closing command.
	Comment_VerbatimLine = C.CXComment_VerbatimLine
	// A full comment attached to a declaration, contains block content.
	Comment_FullComment = C.CXComment_FullComment
)
