package phoenix

// #include "go-clang.h"
import "C"

// A cursor representing some element in the abstract syntax tree for a translation unit. The cursor abstraction unifies the different kinds of entities in a program--declaration, statements, expressions, references to declarations, etc.--under a single "cursor" abstraction with a common set of operations. Common operation for a cursor include: getting the physical location in a source file where the cursor points, getting the name associated with a cursor, and retrieving cursors for any child nodes of a particular cursor. Cursors can be produced in two specific ways. clang_getTranslationUnitCursor() produces a cursor for a translation unit, from which one can use clang_visitChildren() to explore the rest of the translation unit. clang_getCursor() maps from a physical source location to the entity that resides at that location, allowing one to map from the source code into the AST.
type Cursor struct {
	c C.CXCursor
}

// Determine whether two cursors are equivalent.
func EqualCursors(c1, c2 Cursor) bool {
	o := C.clang_equalCursors(c1.c, c2.c)

	return o != C.uint(0)
}

// Returns the Objective-C type encoding for the specified declaration.
func (c Cursor) DeclObjCTypeEncoding() string {
	o := cxstring{C.clang_getDeclObjCTypeEncoding(c.c)}
	defer o.Dispose()

	return o.String()
}

// Returns non-zero if the cursor specifies a Record member that is a bitfield.
func (c Cursor) IsBitField() bool {
	o := C.clang_Cursor_isBitField(c.c)

	return o != C.uint(0)
}

// Returns 1 if the base class specified by the cursor with kind CX_CXXBaseSpecifier is virtual.
func (c Cursor) IsVirtualBase() bool {
	o := C.clang_isVirtualBase(c.c)

	return o != C.uint(0)
}

// Retrieve a Unified Symbol Resolution (USR) for the entity referenced by the given cursor. A Unified Symbol Resolution (USR) is a string that identifies a particular entity (function, class, variable, etc.) within a program. USRs can be compared across translation units to determine, e.g., when references in one translation refer to an entity defined in another translation unit.
func (c Cursor) USR() string {
	o := cxstring{C.clang_getCursorUSR(c.c)}
	defer o.Dispose()

	return o.String()
}

// Retrieve a name for the entity referenced by this cursor.
func (c Cursor) Spelling() string {
	o := cxstring{C.clang_getCursorSpelling(c.c)}
	defer o.Dispose()

	return o.String()
}

// Retrieve the display name for the entity referenced by this cursor. The display name contains extra information that helps identify the cursor, such as the parameters of a function or template or the arguments of a class template specialization.
func (c Cursor) DisplayName() string {
	o := cxstring{C.clang_getCursorDisplayName(c.c)}
	defer o.Dispose()

	return o.String()
}

// Determine whether the declaration pointed to by this cursor is also a definition of that entity.
func (c Cursor) IsCursorDefinition() bool {
	o := C.clang_isCursorDefinition(c.c)

	return o != C.uint(0)
}

// Given a cursor that represents an ObjC method or property declaration, return non-zero if the declaration was affected by "@optional". Returns zero if the cursor is not such a declaration or it is "@required".
func (c Cursor) IsObjCOptional() bool {
	o := C.clang_Cursor_isObjCOptional(c.c)

	return o != C.uint(0)
}

// Returns non-zero if the given cursor is a variadic function or method.
func (c Cursor) IsVariadic() bool {
	o := C.clang_Cursor_isVariadic(c.c)

	return o != C.uint(0)
}

// Given a cursor that represents a declaration, return the associated comment text, including comment markers.
func (c Cursor) RawCommentText() string {
	o := cxstring{C.clang_Cursor_getRawCommentText(c.c)}
	defer o.Dispose()

	return o.String()
}

// Given a cursor that represents a documentable entity (e.g., declaration), return the associated \ paragraph; otherwise return the first paragraph.
func (c Cursor) BriefCommentText() string {
	o := cxstring{C.clang_Cursor_getBriefCommentText(c.c)}
	defer o.Dispose()

	return o.String()
}
