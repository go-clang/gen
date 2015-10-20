package clang

// #include <stdlib.h>
// #include "go-clang.h"
//
import "C"

/**
 * \brief Describes parameter passing direction for \\param or \\arg command.
 */
type CommentParamPassDirection int

const (
	/**
	 * \brief The parameter is an input parameter.
	 */
	CommentParamPassDirection_In = C.CXCommentParamPassDirection_In

	/**
	 * \brief The parameter is an output parameter.
	 */
	CommentParamPassDirection_Out = C.CXCommentParamPassDirection_Out

	/**
	 * \brief The parameter is an input and output parameter.
	 */
	CommentParamPassDirection_InOut = C.CXCommentParamPassDirection_InOut
)
