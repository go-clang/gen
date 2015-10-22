package phoenix

// #include "go-clang.h"
import "C"

// Describes parameter passing direction for \\param or \\arg command.
type CommentParamPassDirection uint32

const (
	// The parameter is an input parameter.
	CommentParamPassDirection_In CommentParamPassDirection = C.CXCommentParamPassDirection_In
	// The parameter is an output parameter.
	CommentParamPassDirection_Out = C.CXCommentParamPassDirection_Out
	// The parameter is an input and output parameter.
	CommentParamPassDirection_InOut = C.CXCommentParamPassDirection_InOut
)
