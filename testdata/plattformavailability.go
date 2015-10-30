package clang

/**
 * \brief A string that describes the platform for which this structure
 * provides availability information.
 *
 * Possible values are "ios" or "macosx".
 */
func (p *PlatformAvailability) Platform() string {
	o := cxstring{p.c.Platform}
	//defer o.Dispose() // done by PlatformAvailability.Dispose()
	return o.String()
}

/**
 * \brief Free the memory associated with a \c CXPlatformAvailability structure.
 */
func (p *PlatformAvailability) Dispose() {
	C.clang_disposeCXPlatformAvailability(&p.c)
}
