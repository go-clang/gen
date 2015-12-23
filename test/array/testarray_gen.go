package clang

// #include "../testdata/test-cases/array.h"
// #include "go-clang.h"
import "C"

type TestArray struct {
	c C.TestArray
}

func (ta TestArray) FunctionWithStructArrayParam(earr EmptyStruct) {
	C.functionWithStructArrayParam(ta.c, C.EmptyStruct(earr))
}

func (ta TestArray) FunctionWithULongArrayParam(larr uint32) {
	C.functionWithULongArrayParam(ta.c, C.ulong(larr))
}

func (ta TestArray) FunctionWithStructArrayParamNoSize(earr EmptyStruct, size int16) {
	C.functionWithStructArrayParamNoSize(ta.c, C.EmptyStruct(earr), C.int(size))
}

func (ta TestArray) FunctionWithULongArrayParamNoSize(larr uint32, size int16) {
	C.functionWithULongArrayParamNoSize(ta.c, C.ulong(larr), C.int(size))
}

func (ta TestArray) Structs() EmptyStruct {
	return EmptyStruct(ta.c.structs)
}

func (ta TestArray) FixedSizedArray() uint32 {
	return uint32(ta.c.fixedSizedArray)
}
