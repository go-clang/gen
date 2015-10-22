package phoenix

// #include "go-clang.h"
import "C"

// Describes a kind of token.
type TokenKind uint32

const (
	// A token that contains some kind of punctuation.
	Token_Punctuation TokenKind = C.CXToken_Punctuation
	// A language keyword.
	Token_Keyword = C.CXToken_Keyword
	// An identifier (that is not a keyword).
	Token_Identifier = C.CXToken_Identifier
	// A numeric, string, or character literal.
	Token_Literal = C.CXToken_Literal
	// A comment.
	Token_Comment = C.CXToken_Comment
)
