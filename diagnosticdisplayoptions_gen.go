package phoenix

// #include "go-clang.h"
import "C"

// Options to control the display of diagnostics. The values in this enum are meant to be combined to customize the behavior of \c clang_formatDiagnostic().
type DiagnosticDisplayOptions int

const (
	/*
	 * \brief Display the source-location information where the
	 * diagnostic was located.
	 *
	 * When set, diagnostics will be prefixed by the file, line, and
	 * (optionally) column to which the diagnostic refers. For example,
	 *
	 * \code
	 * test.c:28: warning: extra tokens at end of #endif directive
	 * \endcode
	 *
	 * This option corresponds to the clang flag \c -fshow-source-location.
	 */
	Diagnostic_DisplaySourceLocation DiagnosticDisplayOptions = C.CXDiagnostic_DisplaySourceLocation
	// If displaying the source-location information of the diagnostic, also include the column number. This option corresponds to the clang flag \c -fshow-column.
	Diagnostic_DisplayColumn = C.CXDiagnostic_DisplayColumn
	// If displaying the source-location information of the diagnostic, also include information about source ranges in a machine-parsable format. This option corresponds to the clang flag \c -fdiagnostics-print-source-range-info.
	Diagnostic_DisplaySourceRanges = C.CXDiagnostic_DisplaySourceRanges
	// Display the option name associated with this diagnostic, if any. The option name displayed (e.g., -Wconversion) will be placed in brackets after the diagnostic text. This option corresponds to the clang flag \c -fdiagnostics-show-option.
	Diagnostic_DisplayOption = C.CXDiagnostic_DisplayOption
	// Display the category number associated with this diagnostic, if any. The category number is displayed within brackets after the diagnostic text. This option corresponds to the clang flag \c -fdiagnostics-show-category=id.
	Diagnostic_DisplayCategoryId = C.CXDiagnostic_DisplayCategoryId
	// Display the category name associated with this diagnostic, if any. The category name is displayed within brackets after the diagnostic text. This option corresponds to the clang flag \c -fdiagnostics-show-category=name.
	Diagnostic_DisplayCategoryName = C.CXDiagnostic_DisplayCategoryName
)
