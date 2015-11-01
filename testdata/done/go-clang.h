#ifndef _GO_CLANG
#define _GO_CLANG 1

#include "clang-c/Index.h"

inline static
CXPlatformAvailability
_goclang_get_platform_availability_at(CXPlatformAvailability* array, int idx) {
  return array[idx];
}

CXPlatformAvailability
_goclang_get_platform_availability_at(CXPlatformAvailability* array, int idx);

#endif /* !_GO_CLANG */
