package phoenix

// #include "go-clang.h"
import "C"

/**
 * \brief Describes the availability of a particular entity, which indicates
 * whether the use of this entity will result in a warning or error due to
 * it being deprecated or unavailable.
 */
type AvailabilityKind int

const (
	/**
	 * \brief The entity is available.
	 */
	Availability_Available AvailabilityKind = C.CXAvailability_Available
	/**
	 * \brief The entity is available, but has been deprecated (and its use is
	 * not recommended).
	 */
	Availability_Deprecated = C.CXAvailability_Deprecated
	/**
	 * \brief The entity is not available; any use of it will be an error.
	 */
	Availability_NotAvailable = C.CXAvailability_NotAvailable
	/**
	 * \brief The entity is available, but not accessible; any use of it will be
	 * an error.
	 */
	Availability_NotAccessible = C.CXAvailability_NotAccessible
)
