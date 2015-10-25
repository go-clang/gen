package phoenix

// #include "go-clang.h"
import "C"

// A single translation unit, which resides in an index.
type TranslationUnit struct {
	c C.CXTranslationUnit
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

// Destroy the specified CXTranslationUnit object.
func (tu TranslationUnit) Dispose() {
	C.clang_disposeTranslationUnit(tu.c)
}

// Return the memory usage of a translation unit. This object should be released with clang_disposeCXTUResourceUsage().
func (tu TranslationUnit) CXTUResourceUsage() TUResourceUsage {
	return TUResourceUsage{C.clang_getCXTUResourceUsage(tu.c)}
}

// Retrieve the cursor that represents the given translation unit. The translation unit cursor can be used to start traversing the various declarations within the given translation unit.
func (tu TranslationUnit) Cursor() Cursor {
	return Cursor{C.clang_getTranslationUnitCursor(tu.c)}
}
