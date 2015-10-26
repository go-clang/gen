package phoenix

// #include "go-clang.h"
import "C"

// A single translation unit, which resides in an index.
type TranslationUnit struct {
	c C.CXTranslationUnit
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
func (tu TranslationUnit) Cursor() Cursor {
	return Cursor{C.clang_getTranslationUnitCursor(tu.c)}
}
