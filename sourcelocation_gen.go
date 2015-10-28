package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// Identifies a specific source location within a translation unit. Use clang_getExpansionLocation() or clang_getSpellingLocation() to map a source location to a particular file, line, and column.
type SourceLocation struct {
	c C.CXSourceLocation
}

func (sl SourceLocation) Ptr_data() []unsafe.Pointer {
	sc := []unsafe.Pointer{}

	length := 2
	goslice := (*[1 << 30]*C.void)(unsafe.Pointer(&sl.c.ptr_data))[:length:length]

	for is := 0; is < length; is++ {
		sc = append(sc, unsafe.Pointer(goslice[is]))
	}

	return sc
}

func (sl SourceLocation) Int_data() uint16 {
	value := uint16(sl.c.int_data)
	return value
}

// Retrieve a NULL (invalid) source location.
func NewNullLocation() SourceLocation {
	return SourceLocation{C.clang_getNullLocation()}
}

// Determine whether two source locations, which must refer into the same translation unit, refer to exactly the same point in the source code. \returns non-zero if the source locations refer to the same location, zero if they refer to different locations.
func (sl SourceLocation) EqualLocations(sl2 SourceLocation) bool {
	o := C.clang_equalLocations(sl.c, sl2.c)

	return o != C.uint(0)
}

// Returns non-zero if the given source location is in a system header.
func (sl SourceLocation) Location_IsInSystemHeader() bool {
	o := C.clang_Location_isInSystemHeader(sl.c)

	return o != C.int(0)
}

// Returns non-zero if the given source location is in the main file of the corresponding translation unit.
func (sl SourceLocation) Location_IsFromMainFile() bool {
	o := C.clang_Location_isFromMainFile(sl.c)

	return o != C.int(0)
}

// Retrieve a source range given the beginning and ending source locations.
func (sl SourceLocation) Range(end SourceLocation) SourceRange {
	return SourceRange{C.clang_getRange(sl.c, end.c)}
}

// Retrieve the file, line, column, and offset represented by the given source location. If the location refers into a macro expansion, retrieves the location of the macro expansion. \param location the location within a source file that will be decomposed into its parts. \param file [out] if non-NULL, will be set to the file to which the given source location points. \param line [out] if non-NULL, will be set to the line to which the given source location points. \param column [out] if non-NULL, will be set to the column to which the given source location points. \param offset [out] if non-NULL, will be set to the offset into the buffer to which the given source location points.
func (sl SourceLocation) ExpansionLocation() (File, uint16, uint16, uint16) {
	var file File
	var line C.uint
	var column C.uint
	var offset C.uint

	C.clang_getExpansionLocation(sl.c, &file.c, &line, &column, &offset)

	return file, uint16(line), uint16(column), uint16(offset)
}

/*
 * \brief Retrieve the file, line, column, and offset represented by
 * the given source location, as specified in a # line directive.
 *
 * Example: given the following source code in a file somefile.c
 *
 * \code
 * #123 "dummy.c" 1
 *
 * static int func(void)
 * {
 *     return 0;
 * }
 * \endcode
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
func (sl SourceLocation) PresumedLocation() (string, uint16, uint16) {
	var filename cxstring
	defer filename.Dispose()
	var line C.uint
	var column C.uint

	C.clang_getPresumedLocation(sl.c, &filename.c, &line, &column)

	return filename.String(), uint16(line), uint16(column)
}

// Legacy API to retrieve the file, line, column, and offset represented by the given source location. This interface has been replaced by the newer interface #clang_getExpansionLocation(). See that interface's documentation for details.
func (sl SourceLocation) InstantiationLocation() (File, uint16, uint16, uint16) {
	var file File
	var line C.uint
	var column C.uint
	var offset C.uint

	C.clang_getInstantiationLocation(sl.c, &file.c, &line, &column, &offset)

	return file, uint16(line), uint16(column), uint16(offset)
}

// Retrieve the file, line, column, and offset represented by the given source location. If the location refers into a macro instantiation, return where the location was originally spelled in the source file. \param location the location within a source file that will be decomposed into its parts. \param file [out] if non-NULL, will be set to the file to which the given source location points. \param line [out] if non-NULL, will be set to the line to which the given source location points. \param column [out] if non-NULL, will be set to the column to which the given source location points. \param offset [out] if non-NULL, will be set to the offset into the buffer to which the given source location points.
func (sl SourceLocation) SpellingLocation() (File, uint16, uint16, uint16) {
	var file File
	var line C.uint
	var column C.uint
	var offset C.uint

	C.clang_getSpellingLocation(sl.c, &file.c, &line, &column, &offset)

	return file, uint16(line), uint16(column), uint16(offset)
}

// Retrieve the file, line, column, and offset represented by the given source location. If the location refers into a macro expansion, return where the macro was expanded or where the macro argument was written, if the location points at a macro argument. \param location the location within a source file that will be decomposed into its parts. \param file [out] if non-NULL, will be set to the file to which the given source location points. \param line [out] if non-NULL, will be set to the line to which the given source location points. \param column [out] if non-NULL, will be set to the column to which the given source location points. \param offset [out] if non-NULL, will be set to the offset into the buffer to which the given source location points.
func (sl SourceLocation) FileLocation() (File, uint16, uint16, uint16) {
	var file File
	var line C.uint
	var column C.uint
	var offset C.uint

	C.clang_getFileLocation(sl.c, &file.c, &line, &column, &offset)

	return file, uint16(line), uint16(column), uint16(offset)
}
