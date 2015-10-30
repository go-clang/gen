package clang

/**
 * \brief Retrieve the unique ID for the given \c file.
 *
 * \param file the file to get the ID for.
 * \param outID stores the returned CXFileUniqueID.
 * \returns If there was a failure getting the unique ID, returns non-zero,
 * otherwise returns 0.
 */
func (f File) GetFileUniqueID() (FileUniqueID, error) {
	var fid FileUniqueID
	o := C.clang_getFileUniqueID(f.c, &fid.c)
	if o != 0 {
		return fid, fmt.Errorf("clang: could not get FileUniqueID (err=%d)", o)
	}
	return fid, nil
}
