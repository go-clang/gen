package phoenix

// #include "go-clang.h"
import "C"

/**
 * \brief Describe the linkage of the entity referred to by a cursor.
 */
type LinkageKind int

const (
	/** \brief This value indicates that no linkage information is available
	 * for a provided CXCursor. */
	Linkage_Invalid LinkageKind = C.CXLinkage_Invalid
	/**
	 * \brief This is the linkage for variables, parameters, and so on that
	 *  have automatic storage.  This covers normal (non-extern) local variables.
	 */
	Linkage_NoLinkage = C.CXLinkage_NoLinkage
	/** \brief This is the linkage for static variables and static functions. */
	Linkage_Internal = C.CXLinkage_Internal
	/** \brief This is the linkage for entities with external linkage that live
	 * in C++ anonymous namespaces.*/
	Linkage_UniqueExternal = C.CXLinkage_UniqueExternal
	/** \brief This is the linkage for entities with true, external linkage. */
	Linkage_External = C.CXLinkage_External
)
