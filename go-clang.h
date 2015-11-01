#ifndef GO_CLANG
#define GO_CLANG

#include <stdlib.h>

#include "clang-c/Index.h"

/* TODO define as constants
#define CINDEX_VERSION_MAJOR 0
#define CINDEX_VERSION_MINOR 20
*/

unsigned go_clang_visit_children(CXCursor c, void *fct);

#endif
