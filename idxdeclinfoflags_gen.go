package phoenix

// #include "go-clang.h"
import "C"

type IdxDeclInfoFlags int

const (
	IdxDeclFlag_Skipped IdxDeclInfoFlags = C.CXIdxDeclFlag_Skipped
)
