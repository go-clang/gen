package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"
import (
	"unsafe"
)

/**
 * \brief A cursor representing some element in the abstract syntax tree for
 * a translation unit.
 *
 * The cursor abstraction unifies the different kinds of entities in a
 * program--declaration, statements, expressions, references to declarations,
 * etc.--under a single "cursor" abstraction with a common set of operations.
 * Common operation for a cursor include: getting the physical location in
 * a source file where the cursor points, getting the name associated with a
 * cursor, and retrieving cursors for any child nodes of a particular cursor.
 *
 * Cursors can be produced in two specific ways.
 * clang_getTranslationUnitCursor() produces a cursor for a translation unit,
 * from which one can use clang_visitChildren() to explore the rest of the
 * translation unit. clang_getCursor() maps from a physical source location
 * to the entity that resides at that location, allowing one to map from the
 * source code into the AST.
 */
type Cursor struct {
	c C.CXCursor
}

/**
 * \brief Determine the availability of the entity that this cursor refers to
 * on any platforms for which availability information is known.
 *
 * \param cursor The cursor to query.
 *
 * \param always_deprecated If non-NULL, will be set to indicate whether the
 * entity is deprecated on all platforms.
 *
 * \param deprecated_message If non-NULL, will be set to the message text
 * provided along with the unconditional deprecation of this entity. The client
 * is responsible for deallocating this string.
 *
 * \param always_unavailable If non-NULL, will be set to indicate whether the
 * entity is unavailable on all platforms.
 *
 * \param unavailable_message If non-NULL, will be set to the message text
 * provided along with the unconditional unavailability of this entity. The
 * client is responsible for deallocating this string.
 *
 * \param availability If non-NULL, an array of CXPlatformAvailability instances
 * that will be populated with platform availability information, up to either
 * the number of platforms for which availability information is available (as
 * returned by this function) or \c availability_size, whichever is smaller.
 *
 * \param availability_size The number of elements available in the
 * \c availability array.
 *
 * \returns The number of platforms (N) for which availability information is
 * available (which is unrelated to \c availability_size).
 *
 * Note that the client is responsible for calling
 * \c clang_disposeCXPlatformAvailability to free each of the
 * platform-availability structures returned. There are
 * \c min(N, availability_size) such structures.
 */
func (c Cursor) PlatformAvailability(availability []PlatformAvailability) (always_deprecated bool, deprecated_msg string, always_unavailable bool, unavailable_msg string) {
	var c_always_deprecated C.int
	var c_deprecated_msg cxstring
	defer c_deprecated_msg.Dispose()
	var c_always_unavailable C.int
	var c_unavailable_msg cxstring
	defer c_unavailable_msg.Dispose()
	c_platforms := make([]C.CXPlatformAvailability, len(availability))

	nn := int(C.clang_getCursorPlatformAvailability(
		c.c,
		&c_always_deprecated,
		&c_deprecated_msg.c,
		&c_always_unavailable,
		&c_unavailable_msg.c,
		&c_platforms[0],
		C.int(len(c_platforms)),
	))

	if c_always_deprecated != 0 {
		always_deprecated = true
	}
	deprecated_msg = c_deprecated_msg.String()

	if c_always_unavailable != 0 {
		always_unavailable = true
	}
	unavailable_msg = c_unavailable_msg.String()

	if nn > len(availability) {
		nn = len(availability)
	}

	availability = make([]PlatformAvailability, nn)
	for i := 0; i < nn; i++ {
		availability[i] = PlatformAvailability{C._goclang_get_platform_availability_at(&c_platforms[0], C.int(i))}
	}

	return
}

// CursorSet is a fast container representing a set of Cursors.
type CursorSet struct {
	c C.CXCursorSet
}

// Contains queries a CursorSet to see if it contains a specific Cursor
func (c CursorSet) Contains(cursor Cursor) bool {
	o := C.clang_CXCursorSet_contains(c.c, cursor.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

// Insert inserts a Cursor into the set and returns false if the cursor was
// already in that set.
func (c CursorSet) Insert(cursor Cursor) bool {
	o := C.clang_CXCursorSet_insert(c.c, cursor.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 * \brief Determine the set of methods that are overridden by the given
 * method.
 *
 * In both Objective-C and C++, a method (aka virtual member function,
 * in C++) can override a virtual method in a base class. For
 * Objective-C, a method is said to override any method in the class's
 * interface (if we're coming from an implementation), its protocols,
 * or its categories, that has the same selector and is of the same
 * kind (class or instance). If no such method exists, the search
 * continues to the class's superclass, its protocols, and its
 * categories, and so on.
 *
 * For C++, a virtual member function overrides any virtual member
 * function with the same signature that occurs in its base
 * classes. With multiple inheritance, a virtual member function can
 * override several virtual member functions coming from different
 * base classes.
 *
 * In all cases, this function determines the immediate overridden
 * method, rather than all of the overridden methods. For example, if
 * a method is originally declared in a class A, then overridden in B
 * (which in inherits from A) and also in C (which inherited from B),
 * then the only overridden method returned from this function when
 * invoked on C's method will be B's method. The client may then
 * invoke this function again, given the previously-found overridden
 * methods, to map out the complete method-override set.
 *
 * \param cursor A cursor representing an Objective-C or C++
 * method. This routine will compute the set of methods that this
 * method overrides.
 *
 * \param overridden A pointer whose pointee will be replaced with a
 * pointer to an array of cursors, representing the set of overridden
 * methods. If there are no overridden methods, the pointee will be
 * set to NULL. The pointee must be freed via a call to
 * \c clang_disposeOverriddenCursors().
 *
 * \param num_overridden A pointer to the number of overridden
 * functions, will be set to the number of overridden functions in the
 * array pointed to by \p overridden.
 */
func (c Cursor) OverriddenCursors() (o OverriddenCursors) {
	C.clang_getOverriddenCursors(c.c, &o.c, &o.n)

	return o
}

type OverriddenCursors struct {
	c *C.CXCursor
	n C.uint
}

// Dispose frees the set of overridden cursors
func (c OverriddenCursors) Dispose() {
	C.clang_disposeOverriddenCursors(c.c)
}

func (c OverriddenCursors) Len() int {
	return int(c.n)
}

func (c OverriddenCursors) At(i int) Cursor {
	if i >= int(c.n) {
		panic("clang: index out of range")
	}
	return Cursor{C._go_clang_ocursor_at(c.c, C.int(i))}
}

/**
 * \brief Retrieve the integer value of an enum constant declaration as a signed
 *  long long.
 *
 * If the cursor does not reference an enum constant declaration, LLONG_MIN is returned.
 * Since this is also potentially a valid constant value, the kind of the cursor
 * must be verified before calling this function.
 */
func (c Cursor) EnumConstantDeclValue() int64 {
	return int64(C.clang_getEnumConstantDeclValue(c.c))
}

/**
 * \brief Retrieve the argument cursor of a function or method.
 *
 * If a cursor that is not a function or method is passed in or the index
 * exceeds the number of arguments, an invalid cursor is returned.
 */
// CINDEX_LINKAGE CXCursor clang_Cursor_getArgument(CXCursor C, unsigned i);
func (c Cursor) Argument(i uint) Cursor {
	o := C.clang_Cursor_getArgument(c.c, C.uint(i))
	return Cursor{o}
}

/**
 * \brief Retrieve a cursor for one of the overloaded declarations referenced
 * by a \c CXCursor_OverloadedDeclRef cursor.
 *
 * \param cursor The cursor whose overloaded declarations are being queried.
 *
 * \param index The zero-based index into the set of overloaded declarations in
 * the cursor.
 *
 * \returns A cursor representing the declaration referenced by the given
 * \c cursor at the specified \c index. If the cursor does not have an
 * associated set of overloaded declarations, or if the index is out of bounds,
 * returns \c clang_getNullCursor();
 */
func (c Cursor) OverloadedDecl(i int) Cursor {
	o := C.clang_getOverloadedDecl(c.c, C.uint(i))
	return Cursor{o}
}

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
	o := C._go_clang_visit_children(c.c, unsafe.Pointer(&visitor))
	if o != C.uint(0) {
		return false
	}
	return true
}

// forceEscapeVisitor is write-only: to force compiler to escape the address
// (else the address can become stale if the goroutine stack needs to grow
// and is forced to move)
// Explained by rsc in https://golang.org/issue/9125
var forceEscapeVisitor *CursorVisitor

//export GoClangCursorVisitor
func GoClangCursorVisitor(cursor, parent C.CXCursor, cfct unsafe.Pointer) (status ChildVisitResult) {
	fct := *(*CursorVisitor)(cfct)
	o := fct(Cursor{cursor}, Cursor{parent})
	return o
}

/**
 * \defgroup CINDEX_CPP C++ AST introspection
 *
 * The routines in this group provide access information in the ASTs specific
 * to C++ language features.
 *
 * @{
 */

/**
 * \brief Given a cursor that references something else, return the source range
 * covering that reference.
 *
 * \param C A cursor pointing to a member reference, a declaration reference, or
 * an operator call.
 * \param NameFlags A bitset with three independent flags:
 * CXNameRange_WantQualifier, CXNameRange_WantTemplateArgs, and
 * CXNameRange_WantSinglePiece.
 * \param PieceIndex For contiguous names or when passing the flag
 * CXNameRange_WantSinglePiece, only one piece with index 0 is
 * available. When the CXNameRange_WantSinglePiece flag is not passed for a
 * non-contiguous names, this index can be used to retreive the individual
 * pieces of the name. See also CXNameRange_WantSinglePiece.
 *
 * \returns The piece of the name pointed to by the given cursor. If there is no
 * name, or if the PieceIndex is out-of-range, a null-cursor will be returned.
 */
func (c Cursor) ReferenceNameRange(flags NameRefFlags, pieceIdx uint) SourceRange {
	o := C.clang_getCursorReferenceNameRange(c.c,
		C.uint(flags), C.uint(pieceIdx))
	return SourceRange{o}
}
