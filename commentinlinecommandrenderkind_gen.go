package phoenix

// #include "go-clang.h"
import "C"

/**
 * \brief The most appropriate rendering mode for an inline command, chosen on
 * command semantics in Doxygen.
 */
type CommentInlineCommandRenderKind int

const (
	/**
	 * \brief Command argument should be rendered in a normal font.
	 */
	CommentInlineCommandRenderKind_Normal CommentInlineCommandRenderKind = C.CXCommentInlineCommandRenderKind_Normal
	/**
	 * \brief Command argument should be rendered in a bold font.
	 */
	CommentInlineCommandRenderKind_Bold = C.CXCommentInlineCommandRenderKind_Bold
	/**
	 * \brief Command argument should be rendered in a monospaced font.
	 */
	CommentInlineCommandRenderKind_Monospaced = C.CXCommentInlineCommandRenderKind_Monospaced
	/**
	 * \brief Command argument should be rendered emphasized (typically italic
	 * font).
	 */
	CommentInlineCommandRenderKind_Emphasized = C.CXCommentInlineCommandRenderKind_Emphasized
)
