package phoenix

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"

// Identifies an array of ranges.
type SourceRangeList struct {
	c C.CXSourceRangeList
}

// The number of ranges in the ranges array.
func (srl SourceRangeList) Count() uint16 {
	return uint16(srl.c.count)
}

// An array of CXSourceRanges.
func (srl SourceRangeList) Ranges() *SourceRange {
	o := srl.c.ranges

	var gop_o *SourceRange
	if o != nil {
		gop_o = &SourceRange{*o}
	}

	return gop_o
}
