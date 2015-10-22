package phoenix

// #include "go-clang.h"
import "C"

// Provides the contents of a file that has not yet been saved to disk. Each CXUnsavedFile instance provides the name of a file on the system along with the current contents of that file that have not yet been saved to disk.
type UnsavedFile struct {
	c C.struct_CXUnsavedFile
}
