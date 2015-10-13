package phoenix

// #include "go-clang.h"
import "C"

/**
 * \brief Categorizes how memory is being used by a translation unit.
 */
type TUResourceUsageKind int

const (
	TUResourceUsage_AST                                TUResourceUsageKind = C.CXTUResourceUsage_AST
	TUResourceUsage_Identifiers                        TUResourceUsageKind = C.CXTUResourceUsage_Identifiers
	TUResourceUsage_Selectors                          TUResourceUsageKind = C.CXTUResourceUsage_Selectors
	TUResourceUsage_GlobalCompletionResults            TUResourceUsageKind = C.CXTUResourceUsage_GlobalCompletionResults
	TUResourceUsage_SourceManagerContentCache          TUResourceUsageKind = C.CXTUResourceUsage_SourceManagerContentCache
	TUResourceUsage_AST_SideTables                     TUResourceUsageKind = C.CXTUResourceUsage_AST_SideTables
	TUResourceUsage_SourceManager_Membuffer_Malloc     TUResourceUsageKind = C.CXTUResourceUsage_SourceManager_Membuffer_Malloc
	TUResourceUsage_SourceManager_Membuffer_MMap       TUResourceUsageKind = C.CXTUResourceUsage_SourceManager_Membuffer_MMap
	TUResourceUsage_ExternalASTSource_Membuffer_Malloc TUResourceUsageKind = C.CXTUResourceUsage_ExternalASTSource_Membuffer_Malloc
	TUResourceUsage_ExternalASTSource_Membuffer_MMap   TUResourceUsageKind = C.CXTUResourceUsage_ExternalASTSource_Membuffer_MMap
	TUResourceUsage_Preprocessor                       TUResourceUsageKind = C.CXTUResourceUsage_Preprocessor
	TUResourceUsage_PreprocessingRecord                TUResourceUsageKind = C.CXTUResourceUsage_PreprocessingRecord
	TUResourceUsage_SourceManager_DataStructures       TUResourceUsageKind = C.CXTUResourceUsage_SourceManager_DataStructures
	TUResourceUsage_Preprocessor_HeaderSearch          TUResourceUsageKind = C.CXTUResourceUsage_Preprocessor_HeaderSearch
	TUResourceUsage_MEMORY_IN_BYTES_BEGIN              TUResourceUsageKind = C.CXTUResourceUsage_MEMORY_IN_BYTES_BEGIN
	TUResourceUsage_MEMORY_IN_BYTES_END                TUResourceUsageKind = C.CXTUResourceUsage_MEMORY_IN_BYTES_END
	TUResourceUsage_First                              TUResourceUsageKind = C.CXTUResourceUsage_First
	TUResourceUsage_Last                               TUResourceUsageKind = C.CXTUResourceUsage_Last
)
