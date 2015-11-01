
 /**
  * \brief Visitor invoked for each file in a translation unit
  *        (used with clang_getInclusions()).
  *
  * This visitor function will be invoked by clang_getInclusions() for each
  * file included (either at the top-level or by \#include directives) within
  * a translation unit.  The first argument is the file being included, and
  * the second and third arguments provide the inclusion stack.  The
  * array is sorted in order of immediate inclusion.  For example,
  * the first element refers to the location that included 'included_file'.
  */
typedef void (*CXInclusionVisitor)(CXFile included_file,
                                   CXSourceLocation* inclusion_stack,
                                   unsigned include_len,
                                   CXClientData client_data);

/**
 * \brief Visit the set of preprocessor inclusions in a translation unit.
 *   The visitor function is called with the provided data for every included
 *   file.  This does not include headers included by the PCH file (unless one
 *   is inspecting the inclusions in the PCH file itself).
 */
void clang_getInclusions(CXTranslationUnit tu,
                                        CXInclusionVisitor visitor,
                                        CXClientData client_data);

const CXIdxObjCContainerDeclInfo *
clang_index_getObjCContainerDeclInfo(const CXIdxDeclInfo *);

const CXIdxObjCInterfaceDeclInfo *
clang_index_getObjCInterfaceDeclInfo(const CXIdxDeclInfo *);

const CXIdxObjCCategoryDeclInfo *
clang_index_getObjCCategoryDeclInfo(const CXIdxDeclInfo *);

const CXIdxObjCProtocolRefListInfo *
clang_index_getObjCProtocolRefListInfo(const CXIdxDeclInfo *);

const CXIdxObjCPropertyDeclInfo *
clang_index_getObjCPropertyDeclInfo(const CXIdxDeclInfo *);

const CXIdxIBOutletCollectionAttrInfo *
clang_index_getIBOutletCollectionAttrInfo(const CXIdxAttrInfo *);

const CXIdxCXXClassDeclInfo *
clang_index_getCXXClassDeclInfo(const CXIdxDeclInfo *);

/**
 * \brief For retrieving a custom CXIdxClientEntity attached to an entity.
 */
CXIdxClientEntity
clang_index_getClientEntity(const CXIdxEntityInfo *);

/**
 * \brief For setting a custom CXIdxClientEntity attached to an entity.
 */
void
clang_index_setClientEntity(const CXIdxEntityInfo *, CXIdxClientEntity);

/**
 * \brief Index the given source file and the translation unit corresponding
 * to that file via callbacks implemented through #IndexerCallbacks.
 *
 * \param client_data pointer data supplied by the client, which will
 * be passed to the invoked callbacks.
 *
 * \param index_callbacks Pointer to indexing callbacks that the client
 * implements.
 *
 * \param index_callbacks_size Size of #IndexerCallbacks structure that gets
 * passed in index_callbacks.
 *
 * \param index_options A bitmask of options that affects how indexing is
 * performed. This should be a bitwise OR of the CXIndexOpt_XXX flags.
 *
 * \param out_TU [out] pointer to store a CXTranslationUnit that can be reused
 * after indexing is finished. Set to NULL if you do not require it.
 *
 * \returns If there is a failure from which the there is no recovery, returns
 * non-zero, otherwise returns 0.
 *
 * The rest of the parameters are the same as #clang_parseTranslationUnit.
 */
int clang_indexSourceFile(CXIndexAction,
                                         CXClientData client_data,
                                         IndexerCallbacks *index_callbacks,
                                         unsigned index_callbacks_size,
                                         unsigned index_options,
                                         const char *source_filename,
                                         const char * const *command_line_args,
                                         int num_command_line_args,
                                         struct CXUnsavedFile *unsaved_files,
                                         unsigned num_unsaved_files,
                                         CXTranslationUnit *out_TU,
                                         unsigned TU_options);

/**
 * \brief Index the given translation unit via callbacks implemented through
 * #IndexerCallbacks.
 *
 * The order of callback invocations is not guaranteed to be the same as
 * when indexing a source file. The high level order will be:
 *
 *   -Preprocessor callbacks invocations
 *   -Declaration/reference callbacks invocations
 *   -Diagnostic callback invocations
 *
 * The parameters are the same as #clang_indexSourceFile.
 *
 * \returns If there is a failure from which the there is no recovery, returns
 * non-zero, otherwise returns 0.
 */
int clang_indexTranslationUnit(CXIndexAction,
                                              CXClientData client_data,
                                              IndexerCallbacks *index_callbacks,
                                              unsigned index_callbacks_size,
                                              unsigned index_options,
                                              CXTranslationUnit);

