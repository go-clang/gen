#ifndef _GO_CLANG
#define _GO_CLANG

#include <stdlib.h>

#include "clang-c/Index.h"

/* TODO define as constants
#define CINDEX_VERSION_MAJOR 0
#define CINDEX_VERSION_MINOR 20
*/

unsigned _go_clang_visit_children(CXCursor c, void *fct);

#endif
