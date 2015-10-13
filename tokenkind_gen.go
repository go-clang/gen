package phoenix

// #include "go-clang.h"
import "C"

/**
 * \brief Describes a kind of token.
 */
type TokenKind int

const (
	/**
	 * \brief A token that contains some kind of punctuation.
	 */
	Token_Punctuation TokenKind = C.CXToken_Punctuation
	/**
	 * \brief A language keyword.
	 */
	Token_Keyword = C.CXToken_Keyword
	/**
	 * \brief An identifier (that is not a keyword).
	 */
	Token_Identifier = C.CXToken_Identifier
	/**
	 * \brief A numeric, string, or character literal.
	 */
	Token_Literal = C.CXToken_Literal
	/**
	 * \brief A comment.
	 */
	Token_Comment = C.CXToken_Comment
)
