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

// Returns non-zero if \p cursor is null.
func (c Cursor) IsNull() bool {
	o := C.clang_Cursor_isNull(c.c)

	return o != C.int(0)
}

// Retrieve the kind of the given cursor.
func (c Cursor) Kind() CursorKind {
	return CursorKind(C.clang_getCursorKind(c.c))
}

// Determine the linkage of the entity referred to by a given cursor.
func (c Cursor) Linkage() LinkageKind {
	return LinkageKind(C.clang_getCursorLinkage(c.c))
}

// Determine the availability of the entity that this cursor refers to, taking the current target platform into account. \param cursor The cursor to query. \returns The availability of the cursor.
func (c Cursor) Availability() AvailabilityKind {
	return AvailabilityKind(C.clang_getCursorAvailability(c.c))
}

// Determine the "language" of the entity referred to by a given cursor.
func (c Cursor) Language() LanguageKind {
	return LanguageKind(C.clang_getCursorLanguage(c.c))
}

// Returns the translation unit that a cursor originated from.
func (c Cursor) TranslationUnit() TranslationUnit {
	return TranslationUnit{C.clang_Cursor_getTranslationUnit(c.c)}
}

/*
 * \brief Determine the semantic parent of the given cursor.
 *
 * The semantic parent of a cursor is the cursor that semantically contains
 * the given \p cursor. For many declarations, the lexical and semantic parents
 * are equivalent (the lexical parent is returned by
 * \c clang_getCursorLexicalParent()). They diverge when declarations or
 * definitions are provided out-of-line. For example:
 *
 * \code
 * class C {
 *  void f();
 * };
 *
 * void C::f() { }
 * \endcode
 *
 * In the out-of-line definition of \c C::f, the semantic parent is the
 * the class \c C, of which this function is a member. The lexical parent is
 * the place where the declaration actually occurs in the source code; in this
 * case, the definition occurs in the translation unit. In general, the
 * lexical parent for a given entity can change without affecting the semantics
 * of the program, and the lexical parent of different declarations of the
 * same entity may be different. Changing the semantic parent of a declaration,
 * on the other hand, can have a major impact on semantics, and redeclarations
 * of a particular entity should all have the same semantic context.
 *
 * In the example above, both declarations of \c C::f have \c C as their
 * semantic context, while the lexical context of the first \c C::f is \c C
 * and the lexical context of the second \c C::f is the translation unit.
 *
 * For global declarations, the semantic parent is the translation unit.
 */
func (c Cursor) SemanticParent() Cursor {
	return Cursor{C.clang_getCursorSemanticParent(c.c)}
}

/*
 * \brief Determine the lexical parent of the given cursor.
 *
 * The lexical parent of a cursor is the cursor in which the given \p cursor
 * was actually written. For many declarations, the lexical and semantic parents
 * are equivalent (the semantic parent is returned by
 * \c clang_getCursorSemanticParent()). They diverge when declarations or
 * definitions are provided out-of-line. For example:
 *
 * \code
 * class C {
 *  void f();
 * };
 *
 * void C::f() { }
 * \endcode
 *
 * In the out-of-line definition of \c C::f, the semantic parent is the
 * the class \c C, of which this function is a member. The lexical parent is
 * the place where the declaration actually occurs in the source code; in this
 * case, the definition occurs in the translation unit. In general, the
 * lexical parent for a given entity can change without affecting the semantics
 * of the program, and the lexical parent of different declarations of the
 * same entity may be different. Changing the semantic parent of a declaration,
 * on the other hand, can have a major impact on semantics, and redeclarations
 * of a particular entity should all have the same semantic context.
 *
 * In the example above, both declarations of \c C::f have \c C as their
 * semantic context, while the lexical context of the first \c C::f is \c C
 * and the lexical context of the second \c C::f is the translation unit.
 *
 * For declarations written in the global scope, the lexical parent is
 * the translation unit.
 */
func (c Cursor) LexicalParent() Cursor {
	return Cursor{C.clang_getCursorLexicalParent(c.c)}
}

// Retrieve the file that is included by the given inclusion directive cursor.
func (c Cursor) IncludedFile() File {
	return File{C.clang_getIncludedFile(c.c)}
}

// Retrieve the physical location of the source constructor referenced by the given cursor. The location of a declaration is typically the location of the name of that declaration, where the name of that declaration would occur if it is unnamed, or some keyword that introduces that particular declaration. The location of a reference is where that reference occurs within the source code.
func (c Cursor) Location() SourceLocation {
	return SourceLocation{C.clang_getCursorLocation(c.c)}
}

// Retrieve the physical extent of the source construct referenced by the given cursor. The extent of a cursor starts with the file/line/column pointing at the first character within the source construct that the cursor refers to and ends with the last character withinin that source construct. For a declaration, the extent covers the declaration itself. For a reference, the extent covers the location of the reference (e.g., where the referenced entity was actually used).
func (c Cursor) Extent() SourceRange {
	return SourceRange{C.clang_getCursorExtent(c.c)}
}

// Retrieve the type of a CXCursor (if any).
func (c Cursor) Type() Type {
	return Type{C.clang_getCursorType(c.c)}
}

// Retrieve the underlying type of a typedef declaration. If the cursor does not reference a typedef declaration, an invalid type is returned.
func (c Cursor) TypedefDeclUnderlyingType() Type {
	return Type{C.clang_getTypedefDeclUnderlyingType(c.c)}
}

// Retrieve the integer type of an enum declaration. If the cursor does not reference an enum declaration, an invalid type is returned.
func (c Cursor) EnumDeclIntegerType() Type {
	return Type{C.clang_getEnumDeclIntegerType(c.c)}
}

// Returns the Objective-C type encoding for the specified declaration.
func (c Cursor) DeclObjCTypeEncoding() string {
	o := cxstring{C.clang_getDeclObjCTypeEncoding(c.c)}
	defer o.Dispose()

	return o.String()
}

// Retrieve the result type associated with a given cursor. This only returns a valid type if the cursor refers to a function or method.
func (c Cursor) ResultType() Type {
	return Type{C.clang_getCursorResultType(c.c)}
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

// Returns the access control level for the referenced object. If the cursor refers to a C++ declaration, its access control level within its parent scope is returned. Otherwise, if the cursor refers to a base specifier or access specifier, the specifier itself is returned.
func (c Cursor) CXXAccessSpecifier() AccessSpecifier {
	return AccessSpecifier(C.clang_getCXXAccessSpecifier(c.c))
}

// For cursors representing an iboutletcollection attribute, this function returns the collection element type.
func (c Cursor) IBOutletCollectionType() Type {
	return Type{C.clang_getIBOutletCollectionType(c.c)}
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

// For a cursor that is a reference, retrieve a cursor representing the entity that it references. Reference cursors refer to other entities in the AST. For example, an Objective-C superclass reference cursor refers to an Objective-C class. This function produces the cursor for the Objective-C class from the cursor for the superclass reference. If the input cursor is a declaration or definition, it returns that declaration or definition unchanged. Otherwise, returns the NULL cursor.
func (c Cursor) Referenced() Cursor {
	return Cursor{C.clang_getCursorReferenced(c.c)}
}

/*
 *  \brief For a cursor that is either a reference to or a declaration
 *  of some entity, retrieve a cursor that describes the definition of
 *  that entity.
 *
 *  Some entities can be declared multiple times within a translation
 *  unit, but only one of those declarations can also be a
 *  definition. For example, given:
 *
 *  \code
 *  int f(int, int);
 *  int g(int x, int y) { return f(x, y); }
 *  int f(int a, int b) { return a + b; }
 *  int f(int, int);
 *  \endcode
 *
 *  there are three declarations of the function "f", but only the
 *  second one is a definition. The clang_getCursorDefinition()
 *  function will take any cursor pointing to a declaration of "f"
 *  (the first or fourth lines of the example) or a cursor referenced
 *  that uses "f" (the call to "f' inside "g") and will return a
 *  declaration cursor pointing to the definition (the second "f"
 *  declaration).
 *
 *  If given a cursor for which there is no corresponding definition,
 *  e.g., because there is no definition of that entity within this
 *  translation unit, returns a NULL cursor.
 */
func (c Cursor) Definition() Cursor {
	return Cursor{C.clang_getCursorDefinition(c.c)}
}

// Determine whether the declaration pointed to by this cursor is also a definition of that entity.
func (c Cursor) IsCursorDefinition() bool {
	o := C.clang_isCursorDefinition(c.c)

	return o != C.uint(0)
}

/*
 * \brief Retrieve the canonical cursor corresponding to the given cursor.
 *
 * In the C family of languages, many kinds of entities can be declared several
 * times within a single translation unit. For example, a structure type can
 * be forward-declared (possibly multiple times) and later defined:
 *
 * \code
 * struct X;
 * struct X;
 * struct X {
 *   int member;
 * };
 * \endcode
 *
 * The declarations and the definition of \c X are represented by three
 * different cursors, all of which are declarations of the same underlying
 * entity. One of these cursor is considered the "canonical" cursor, which
 * is effectively the representative for the underlying entity. One can
 * determine if two cursors are declarations of the same underlying entity by
 * comparing their canonical cursors.
 *
 * \returns The canonical cursor for the entity referred to by the given cursor.
 */
func (c Cursor) CanonicalCursor() Cursor {
	return Cursor{C.clang_getCanonicalCursor(c.c)}
}

// Given a cursor pointing to a C++ method call or an ObjC message, returns non-zero if the method/message is "dynamic", meaning: For a C++ method: the call is virtual. For an ObjC message: the receiver is an object instance, not 'super' or a specific class. If the method/message is "static" or the cursor does not point to a method/message, it will return zero.
func (c Cursor) IsDynamicCall() bool {
	o := C.clang_Cursor_isDynamicCall(c.c)

	return o != C.int(0)
}

// Given a cursor pointing to an ObjC message, returns the CXType of the receiver.
func (c Cursor) ReceiverType() Type {
	return Type{C.clang_Cursor_getReceiverType(c.c)}
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

// Given a cursor that represents a declaration, return the associated comment's source range. The range may include multiple consecutive comments with whitespace in between.
func (c Cursor) CommentRange() SourceRange {
	return SourceRange{C.clang_Cursor_getCommentRange(c.c)}
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

// Given a cursor that represents a documentable entity (e.g., declaration), return the associated parsed comment as a \c CXComment_FullComment AST node.
func (c Cursor) ParsedComment() Comment {
	return Comment{C.clang_Cursor_getParsedComment(c.c)}
}

// Given a CXCursor_ModuleImportDecl cursor, return the associated module.
func (c Cursor) Module() Module {
	return Module{C.clang_Cursor_getModule(c.c)}
}

// Given a cursor that represents a template, determine the cursor kind of the specializations would be generated by instantiating the template. This routine can be used to determine what flavor of function template, class template, or class template partial specialization is stored in the cursor. For example, it can describe whether a class template cursor is declared with "struct", "class" or "union". \param C The cursor to query. This cursor should represent a template declaration. \returns The cursor kind of the specializations that would be generated by instantiating the template \p C. If \p C is not a template, returns \c CXCursor_NoDeclFound.
func (c Cursor) TemplateCursorKind() CursorKind {
	return CursorKind(C.clang_getTemplateCursorKind(c.c))
}

// Given a cursor that may represent a specialization or instantiation of a template, retrieve the cursor that represents the template that it specializes or from which it was instantiated. This routine determines the template involved both for explicit specializations of templates and for implicit instantiations of the template, both of which are referred to as "specializations". For a class template specialization (e.g., \c std::vector<bool>), this routine will return either the primary template (\c std::vector) or, if the specialization was instantiated from a class template partial specialization, the class template partial specialization. For a class template partial specialization and a function template specialization (including instantiations), this this routine will return the specialized template. For members of a class template (e.g., member functions, member classes, or static data members), returns the specialized or instantiated member. Although not strictly "templates" in the C++ language, members of class templates have the same notions of specializations and instantiations that templates do, so this routine treats them similarly. \param C A cursor that may be a specialization of a template or a member of a template. \returns If the given cursor is a specialization or instantiation of a template or a member thereof, the template or member that it specializes or from which it was instantiated. Otherwise, returns a NULL cursor.
func (c Cursor) SpecializedCursorTemplate() Cursor {
	return Cursor{C.clang_getSpecializedCursorTemplate(c.c)}
}

// Retrieve a completion string for an arbitrary declaration or macro definition cursor. \param cursor The cursor to query. \returns A non-context-sensitive completion string for declaration and macro definition cursors, or NULL for other kinds of cursors.
func (c Cursor) CompletionString() CompletionString {
	return CompletionString{C.clang_getCursorCompletionString(c.c)}
}
