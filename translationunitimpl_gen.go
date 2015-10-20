package phoenix

// #include "go-clang.h"
import "C"

type TranslationUnitImpl struct {
	c C.CXTranslationUnitImpl
}
