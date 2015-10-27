package phoenix

// #include "go-clang.h"
import "C"

// A single translation unit, which resides in an index.
type TranslationUnit struct {
	c C.CXTranslationUnit
}

// Retrieves the source location associated with a given file/line/column in a particular translation unit.
func (tu TranslationUnit) Location(file File, line uint16, column uint16) SourceLocation {
	return SourceLocation{C.clang_getLocation(tu.c, file.c, C.uint(line), C.uint(column))}
}

// Retrieves the source location associated with a given character offset in a particular translation unit.
func (tu TranslationUnit) LocationForOffset(file File, offset uint16) SourceLocation {
	return SourceLocation{C.clang_getLocationForOffset(tu.c, file.c, C.uint(offset))}
}

// Determine the number of diagnostics produced for the given translation unit.
func (tu TranslationUnit) NumDiagnostics() uint16 {
	return uint16(C.clang_getNumDiagnostics(tu.c))
}

// Retrieve the complete set of diagnostics associated with a translation unit. \param Unit the translation unit to query.
func (tu TranslationUnit) DiagnosticSetFromTU() DiagnosticSet {
	return DiagnosticSet{C.clang_getDiagnosticSetFromTU(tu.c)}
}

// Get the original translation unit source file name.
func (tu TranslationUnit) Spelling() string {
	o := cxstring{C.clang_getTranslationUnitSpelling(tu.c)}
	defer o.Dispose()

	return o.String()
}

// Returns the set of flags that is suitable for saving a translation unit. The set of flags returned provide options for \c clang_saveTranslationUnit() by default. The returned flag set contains an unspecified set of options that save translation units with the most commonly-requested data.
func (tu TranslationUnit) DefaultSaveOptions() uint16 {
	return uint16(C.clang_defaultSaveOptions(tu.c))
}

// Destroy the specified CXTranslationUnit object.
func (tu TranslationUnit) Dispose() {
	C.clang_disposeTranslationUnit(tu.c)
}

// Returns the set of flags that is suitable for reparsing a translation unit. The set of flags returned provide options for \c clang_reparseTranslationUnit() by default. The returned flag set contains an unspecified set of optimizations geared toward common uses of reparsing. The set of optimizations enabled may change from one version to the next.
func (tu TranslationUnit) DefaultReparseOptions() uint16 {
	return uint16(C.clang_defaultReparseOptions(tu.c))
}

// Return the memory usage of a translation unit. This object should be released with clang_disposeCXTUResourceUsage().
func (tu TranslationUnit) TUResourceUsage() TUResourceUsage {
	return TUResourceUsage{C.clang_getCXTUResourceUsage(tu.c)}
}

// Retrieve the cursor that represents the given translation unit. The translation unit cursor can be used to start traversing the various declarations within the given translation unit.
func (tu TranslationUnit) TranslationUnitCursor() Cursor {
	return Cursor{C.clang_getTranslationUnitCursor(tu.c)}
}

// Map a source location to the cursor that describes the entity at that location in the source code. clang_getCursor() maps an arbitrary source location within a translation unit down to the most specific cursor that describes the entity at that location. For example, given an expression \c x + y, invoking clang_getCursor() with a source location pointing to "x" will return the cursor for "x"; similarly for "y". If the cursor points anywhere between "x" or "y" (e.g., on the + or the whitespace around it), clang_getCursor() will return a cursor referring to the "+" expression. \returns a cursor representing the entity at the given source location, or a NULL cursor if no such entity can be found.
func (tu TranslationUnit) Cursor(sl SourceLocation) Cursor {
	return Cursor{C.clang_getCursor(tu.c, sl.c)}
}

// Retrieve the source location of the given token.
func (tu TranslationUnit) TokenLocation(t Token) SourceLocation {
	return SourceLocation{C.clang_getTokenLocation(tu.c, t.c)}
}

// Retrieve a source range that covers the given token.
func (tu TranslationUnit) TokenExtent(t Token) SourceRange {
	return SourceRange{C.clang_getTokenExtent(tu.c, t.c)}
}

// Find #import/#include directives in a specific file. \param TU translation unit containing the file to query. \param file to search for #import/#include directives. \param visitor callback that will receive pairs of CXCursor/CXSourceRange for each directive found. \returns one of the CXResult enumerators.
func (tu TranslationUnit) FindIncludesInFile(file File, visitor CursorAndRangeVisitor) Result {
	return Result(C.clang_findIncludesInFile(tu.c, file.c, visitor.c))
}
