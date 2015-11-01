package clang

/**
 * \brief Free the memory associated with a \c CXPlatformAvailability structure.
 */
func (p *PlatformAvailability) Dispose() {
	C.clang_disposeCXPlatformAvailability(&p.c)
}
