package phoenix

// #include "go-clang.h"
import "C"

// Describes how the traversal of the children of a particular cursor should proceed after visiting a particular child cursor. A value of this enumeration type should be returned by each \c CXCursorVisitor to indicate how clang_visitChildren() proceed.
type ChildVisitResult int

const (
	// Terminates the cursor traversal.
	ChildVisit_Break ChildVisitResult = C.CXChildVisit_Break
	// Continues the cursor traversal with the next sibling of the cursor just visited, without visiting its children.
	ChildVisit_Continue = C.CXChildVisit_Continue
	// Recursively traverse the children of this cursor, using the same visitor and client data.
	ChildVisit_Recurse = C.CXChildVisit_Recurse
)
