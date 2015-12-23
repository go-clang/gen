package clang

// #include "../testdata/test-cases/array.h"
// #include "go-clang.h"
import "C"

type EmptyStruct struct {
	c C.EmptyStruct
}
