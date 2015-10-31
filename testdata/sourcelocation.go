package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"

// ExpansionLocation returns the file, line, column, and offset represented by
// the given source location.
//
// If the location refers into a macro expansion, retrieves the
// location of the macro expansion.
//
// file: if non-NULL, will be set to the file to which the given
// source location points.
//
// line: if non-NULL, will be set to the line to which the given
// source location points.
//
// column: if non-NULL, will be set to the column to which the given
// source location points.
//
// offset: if non-NULL, will be set to the offset into the
// buffer to which the given source location points.
func (l SourceLocation) ExpansionLocation() (f File, line, column, offset uint) {
	cline := C.uint(0)
	ccol := C.uint(0)
	coff := C.uint(0)
	// FIXME: undefined reference to `clang_getExpansionLocation'
	C.clang_getInstantiationLocation(l.c, &f.c, &cline, &ccol, &coff)
	line = uint(cline)
	column = uint(ccol)
	offset = uint(coff)

	return
}

/**
 * \brief Retrieve the file, line, column, and offset represented by
 * the given source location, as specified in a # line directive.
 *
 * Example: given the following source code in a file somefile.c
 *
 * #123 "dummy.c" 1
 *
 * static int func(void)
 * {
 *     return 0;
 * }
 *
 * the location information returned by this function would be
 *
 * File: dummy.c Line: 124 Column: 12
 *
 * whereas clang_getExpansionLocation would have returned
 *
 * File: somefile.c Line: 3 Column: 12
 *
 * \param location the location within a source file that will be decomposed
 * into its parts.
 *
 * \param filename [out] if non-NULL, will be set to the filename of the
 * source location. Note that filenames returned will be for "virtual" files,
 * which don't necessarily exist on the machine running clang - e.g. when
 * parsing preprocessed output obtained from a different environment. If
 * a non-NULL value is passed in, remember to dispose of the returned value
 * using \c clang_disposeString() once you've finished with it. For an invalid
 * source location, an empty string is returned.
 *
 * \param line [out] if non-NULL, will be set to the line number of the
 * source location. For an invalid source location, zero is returned.
 *
 * \param column [out] if non-NULL, will be set to the column number of the
 * source location. For an invalid source location, zero is returned.
 */
func (l SourceLocation) PresumedLocation() (fname string, line, column uint) {

	cname := cxstring{}
	defer cname.Dispose()
	cline := C.uint(0)
	ccol := C.uint(0)
	C.clang_getPresumedLocation(l.c, &cname.c, &cline, &ccol)
	fname = cname.String()
	line = uint(cline)
	column = uint(ccol)
	return
}
