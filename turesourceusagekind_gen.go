package phoenix

// #include "go-clang.h"
import "C"

// Categorizes how memory is being used by a translation unit.
type TUResourceUsageKind int

const (
	//
	TUResourceUsage_AST TUResourceUsageKind = C.CXTUResourceUsage_AST
	//
	TUResourceUsage_Identifiers = C.CXTUResourceUsage_Identifiers
	//
	TUResourceUsage_Selectors = C.CXTUResourceUsage_Selectors
	//
	TUResourceUsage_GlobalCompletionResults = C.CXTUResourceUsage_GlobalCompletionResults
	//
	TUResourceUsage_SourceManagerContentCache = C.CXTUResourceUsage_SourceManagerContentCache
	//
	TUResourceUsage_AST_SideTables = C.CXTUResourceUsage_AST_SideTables
	//
	TUResourceUsage_SourceManager_Membuffer_Malloc = C.CXTUResourceUsage_SourceManager_Membuffer_Malloc
	//
	TUResourceUsage_SourceManager_Membuffer_MMap = C.CXTUResourceUsage_SourceManager_Membuffer_MMap
	//
	TUResourceUsage_ExternalASTSource_Membuffer_Malloc = C.CXTUResourceUsage_ExternalASTSource_Membuffer_Malloc
	//
	TUResourceUsage_ExternalASTSource_Membuffer_MMap = C.CXTUResourceUsage_ExternalASTSource_Membuffer_MMap
	//
	TUResourceUsage_Preprocessor = C.CXTUResourceUsage_Preprocessor
	//
	TUResourceUsage_PreprocessingRecord = C.CXTUResourceUsage_PreprocessingRecord
	//
	TUResourceUsage_SourceManager_DataStructures = C.CXTUResourceUsage_SourceManager_DataStructures
	//
	TUResourceUsage_Preprocessor_HeaderSearch = C.CXTUResourceUsage_Preprocessor_HeaderSearch
	//
	TUResourceUsage_MEMORY_IN_BYTES_BEGIN = C.CXTUResourceUsage_MEMORY_IN_BYTES_BEGIN
	//
	TUResourceUsage_MEMORY_IN_BYTES_END = C.CXTUResourceUsage_MEMORY_IN_BYTES_END
	//
	TUResourceUsage_First = C.CXTUResourceUsage_First
	//
	TUResourceUsage_Last = C.CXTUResourceUsage_Last
)
