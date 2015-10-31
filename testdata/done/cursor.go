package clang

/**
 * \brief Determine the availability of the entity that this cursor refers to
 * on any platforms for which availability information is known.
 *
 * \param cursor The cursor to query.
 *
 * \param always_deprecated If non-NULL, will be set to indicate whether the
 * entity is deprecated on all platforms.
 *
 * \param deprecated_message If non-NULL, will be set to the message text
 * provided along with the unconditional deprecation of this entity. The client
 * is responsible for deallocating this string.
 *
 * \param always_unavailable If non-NULL, will be set to indicate whether the
 * entity is unavailable on all platforms.
 *
 * \param unavailable_message If non-NULL, will be set to the message text
 * provided along with the unconditional unavailability of this entity. The
 * client is responsible for deallocating this string.
 *
 * \param availability If non-NULL, an array of CXPlatformAvailability instances
 * that will be populated with platform availability information, up to either
 * the number of platforms for which availability information is available (as
 * returned by this function) or \c availability_size, whichever is smaller.
 *
 * \param availability_size The number of elements available in the
 * \c availability array.
 *
 * \returns The number of platforms (N) for which availability information is
 * available (which is unrelated to \c availability_size).
 *
 * Note that the client is responsible for calling
 * \c clang_disposeCXPlatformAvailability to free each of the
 * platform-availability structures returned. There are
 * \c min(N, availability_size) such structures.
 */
func (c Cursor) PlatformAvailability(availability []PlatformAvailability) (always_deprecated bool, deprecated_msg string, always_unavailable bool, unavailable_msg string) {
	var c_always_deprecated C.int
	var c_deprecated_msg cxstring
	defer c_deprecated_msg.Dispose()
	var c_always_unavailable C.int
	var c_unavailable_msg cxstring
	defer c_unavailable_msg.Dispose()
	c_platforms := make([]C.CXPlatformAvailability, len(availability))

	nn := int(C.clang_getCursorPlatformAvailability(
		c.c,
		&c_always_deprecated,
		&c_deprecated_msg.c,
		&c_always_unavailable,
		&c_unavailable_msg.c,
		&c_platforms[0],
		C.int(len(c_platforms)),
	))

	if c_always_deprecated != 0 {
		always_deprecated = true
	}
	deprecated_msg = c_deprecated_msg.String()

	if c_always_unavailable != 0 {
		always_unavailable = true
	}
	unavailable_msg = c_unavailable_msg.String()

	if nn > len(availability) {
		nn = len(availability)
	}

	availability = make([]PlatformAvailability, nn)
	for i := 0; i < nn; i++ {
		availability[i] = PlatformAvailability{C._goclang_get_platform_availability_at(&c_platforms[0], C.int(i))}
	}

	return
}

// CursorSet is a fast container representing a set of Cursors.
type CursorSet struct {
	c C.CXCursorSet
}

/**
 * \brief Determine the set of methods that are overridden by the given
 * method.
 *
 * In both Objective-C and C++, a method (aka virtual member function,
 * in C++) can override a virtual method in a base class. For
 * Objective-C, a method is said to override any method in the class's
 * interface (if we're coming from an implementation), its protocols,
 * or its categories, that has the same selector and is of the same
 * kind (class or instance). If no such method exists, the search
 * continues to the class's superclass, its protocols, and its
 * categories, and so on.
 *
 * For C++, a virtual member function overrides any virtual member
 * function with the same signature that occurs in its base
 * classes. With multiple inheritance, a virtual member function can
 * override several virtual member functions coming from different
 * base classes.
 *
 * In all cases, this function determines the immediate overridden
 * method, rather than all of the overridden methods. For example, if
 * a method is originally declared in a class A, then overridden in B
 * (which in inherits from A) and also in C (which inherited from B),
 * then the only overridden method returned from this function when
 * invoked on C's method will be B's method. The client may then
 * invoke this function again, given the previously-found overridden
 * methods, to map out the complete method-override set.
 *
 * \param cursor A cursor representing an Objective-C or C++
 * method. This routine will compute the set of methods that this
 * method overrides.
 *
 * \param overridden A pointer whose pointee will be replaced with a
 * pointer to an array of cursors, representing the set of overridden
 * methods. If there are no overridden methods, the pointee will be
 * set to NULL. The pointee must be freed via a call to
 * \c clang_disposeOverriddenCursors().
 *
 * \param num_overridden A pointer to the number of overridden
 * functions, will be set to the number of overridden functions in the
 * array pointed to by \p overridden.
 */
func (c Cursor) OverriddenCursors() (o OverriddenCursors) {
	C.clang_getOverriddenCursors(c.c, &o.c, &o.n)

	return o
}

type OverriddenCursors struct {
	c *C.CXCursor
	n C.uint
}

// Dispose frees the set of overridden cursors
func (c OverriddenCursors) Dispose() {
	C.clang_disposeOverriddenCursors(c.c)
}

func (c OverriddenCursors) Len() int {
	return int(c.n)
}

func (c OverriddenCursors) At(i int) Cursor {
	if i >= int(c.n) {
		panic("clang: index out of range")
	}
	return Cursor{C._go_clang_ocursor_at(c.c, C.int(i))}
}
