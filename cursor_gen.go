package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// A cursor representing some element in the abstract syntax tree for a translation unit. The cursor abstraction unifies the different kinds of entities in a program--declaration, statements, expressions, references to declarations, etc.--under a single "cursor" abstraction with a common set of operations. Common operation for a cursor include: getting the physical location in a source file where the cursor points, getting the name associated with a cursor, and retrieving cursors for any child nodes of a particular cursor. Cursors can be produced in two specific ways. clang_getTranslationUnitCursor() produces a cursor for a translation unit, from which one can use clang_visitChildren() to explore the rest of the translation unit. clang_getCursor() maps from a physical source location to the entity that resides at that location, allowing one to map from the source code into the AST.
type Cursor struct {
	c C.CXCursor
}

func (c Cursor) Kind() CursorKind {
	value := CursorKind(c.c.kind)
	return value
}

func (c Cursor) Xdata() int32 {
	value := int32(c.c.xdata)
	return value
}

func (c Cursor) Data() []unsafe.Pointer {
	sc := []unsafe.Pointer{}

	length := 3
	goslice := (*[1 << 30]*C.void)(unsafe.Pointer(&c.c.data))[:length:length]

	for is := 0; is < length; is++ {
		sc = append(sc, unsafe.Pointer(goslice[is]))
	}

	return sc
}

// Determine whether two cursors are equivalent.
func (c Cursor) EqualCursors(c2 Cursor) bool {
	o := C.clang_equalCursors(c.c, c2.c)

	return o != C.uint(0)
}

// Retrieve the argument cursor of a function or method. The argument cursor can be determined for calls as well as for declarations of functions or methods. For other cursors and for invalid indices, an invalid cursor is returned.
func (c Cursor) Argument(i uint16) Cursor {
	return Cursor{C.clang_Cursor_getArgument(c.c, C.uint(i))}
}

// Retrieve a cursor for one of the overloaded declarations referenced by a \c CXCursor_OverloadedDeclRef cursor. \param cursor The cursor whose overloaded declarations are being queried. \param index The zero-based index into the set of overloaded declarations in the cursor. \returns A cursor representing the declaration referenced by the given \c cursor at the specified \c index. If the cursor does not have an associated set of overloaded declarations, or if the index is out of bounds, returns \c clang_getNullCursor();
func (c Cursor) OverloadedDecl(index uint16) Cursor {
	return Cursor{C.clang_getOverloadedDecl(c.c, C.uint(index))}
}

// Retrieve a range for a piece that forms the cursors spelling name. Most of the times there is only one range for the complete spelling but for objc methods and objc message expressions, there are multiple pieces for each selector identifier. \param pieceIndex the index of the spelling name piece. If this is greater than the actual number of pieces, it will return a NULL (invalid) range. \param options Reserved.
func (c Cursor) SpellingNameRange(pieceIndex uint16, options uint16) SourceRange {
	return SourceRange{C.clang_Cursor_getSpellingNameRange(c.c, C.uint(pieceIndex), C.uint(options))}
}

// Given a cursor that represents a property declaration, return the associated property attributes. The bits are formed from \c CXObjCPropertyAttrKind. \param reserved Reserved for future use, pass 0.
func (c Cursor) ObjCPropertyAttributes(reserved uint16) uint16 {
	return uint16(C.clang_Cursor_getObjCPropertyAttributes(c.c, C.uint(reserved)))
}

// Given a cursor that references something else, return the source range covering that reference. \param C A cursor pointing to a member reference, a declaration reference, or an operator call. \param NameFlags A bitset with three independent flags: CXNameRange_WantQualifier, CXNameRange_WantTemplateArgs, and CXNameRange_WantSinglePiece. \param PieceIndex For contiguous names or when passing the flag CXNameRange_WantSinglePiece, only one piece with index 0 is available. When the CXNameRange_WantSinglePiece flag is not passed for a non-contiguous names, this index can be used to retrieve the individual pieces of the name. See also CXNameRange_WantSinglePiece. \returns The piece of the name pointed to by the given cursor. If there is no name, or if the PieceIndex is out-of-range, a null-cursor will be returned.
func (c Cursor) ReferenceNameRange(NameFlags uint16, PieceIndex uint16) SourceRange {
	return SourceRange{C.clang_getCursorReferenceNameRange(c.c, C.uint(NameFlags), C.uint(PieceIndex))}
}

// Find references of a declaration in a specific file. \param cursor pointing to a declaration or a reference of one. \param file to search for references. \param visitor callback that will receive pairs of CXCursor/CXSourceRange for each reference found. The CXSourceRange will point inside the file; if the reference is inside a macro (and not a macro argument) the CXSourceRange will be invalid. \returns one of the CXResult enumerators.
func (c Cursor) FindReferencesInFile(file File, visitor CursorAndRangeVisitor) Result {
	return Result(C.clang_findReferencesInFile(c.c, file.c, visitor.c))
}
