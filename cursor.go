package phoenix

// #include "go-clang.h"
import "C"
import (
	"unsafe"
)

// CursorVisitor does the following.
/**
 * \brief Visitor invoked for each cursor found by a traversal.
 *
 * This visitor function will be invoked for each cursor found by
 * clang_visitCursorChildren(). Its first argument is the cursor being
 * visited, its second argument is the parent visitor for that cursor,
 * and its third argument is the client data provided to
 * clang_visitCursorChildren().
 *
 * The visitor should return one of the \c CXChildVisitResult values
 * to direct clang_visitCursorChildren().
 */
type CursorVisitor func(cursor, parent Cursor) (status ChildVisitResult)

// GoClangCursorVisitor calls the cursor visitor
//export GoClangCursorVisitor
func GoClangCursorVisitor(cursor C.CXCursor, parent C.CXCursor, cfct unsafe.Pointer) (status ChildVisitResult) {
	fct := *(*CursorVisitor)(cfct)

	return fct(Cursor{cursor}, Cursor{parent})
}

// forceEscapeVisitor is write-only: to force compiler to escape the address
// (else the address can become stale if the goroutine stack needs to grow
// and is forced to move)
// Explained by rsc in https://golang.org/issue/9125
var forceEscapeVisitor *CursorVisitor

// Visit does the following.
/**
 * \brief Visit the children of a particular cursor.
 *
 * This function visits all the direct children of the given cursor,
 * invoking the given \p visitor function with the cursors of each
 * visited child. The traversal may be recursive, if the visitor returns
 * \c CXChildVisit_Recurse. The traversal may also be ended prematurely, if
 * the visitor returns \c CXChildVisit_Break.
 *
 * \param parent the cursor whose child may be visited. All kinds of
 * cursors can be visited, including invalid cursors (which, by
 * definition, have no children).
 *
 * \param visitor the visitor function that will be invoked for each
 * child of \p parent.
 *
 * \param client_data pointer data supplied by the client, which will
 * be passed to the visitor each time it is invoked.
 *
 * \returns a non-zero value if the traversal was terminated
 * prematurely by the visitor returning \c CXChildVisit_Break.
 */
func (c Cursor) Visit(visitor CursorVisitor) bool {
	forceEscapeVisitor = &visitor

	o := C.go_clang_visit_children(c.c, unsafe.Pointer(&visitor))

	return o == C.uint(0)
}
