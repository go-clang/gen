package phoenix

// #include "go-clang.h"
import "C"

// Describes the kind of error that occurred (if any) in a call to \c clang_saveTranslationUnit().
type SaveError int

const (
	// Indicates that no error occurred while saving a translation unit.
	SaveError_None SaveError = C.CXSaveError_None
	// Indicates that an unknown error occurred while attempting to save the file. This error typically indicates that file I/O failed when attempting to write the file.
	SaveError_Unknown = C.CXSaveError_Unknown
	// Indicates that errors during translation prevented this attempt to save the translation unit. Errors that prevent the translation unit from being saved can be extracted using \c clang_getNumDiagnostics() and \c clang_getDiagnostic().
	SaveError_TranslationErrors = C.CXSaveError_TranslationErrors
	// Indicates that the translation unit to be saved was somehow invalid (e.g., NULL).
	SaveError_InvalidTU = C.CXSaveError_InvalidTU
)
