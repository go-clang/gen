package phoenix

// #include "go-clang.h"
import "C"

/**
 * \brief Describes the kind of entity that a cursor refers to.
 */
type CursorKind int

const (
	/**
	 * \brief A declaration whose specific kind is not exposed via this
	 * interface.
	 *
	 * Unexposed declarations have the same operations as any other kind
	 * of declaration; one can extract their location information,
	 * spelling, find their definitions, etc. However, the specific kind
	 * of the declaration is not reported.
	 */
	Cursor_UnexposedDecl CursorKind = C.CXCursor_UnexposedDecl
	/** \brief A C or C++ struct. */
	Cursor_StructDecl CursorKind = C.CXCursor_StructDecl
	/** \brief A C or C++ union. */
	Cursor_UnionDecl CursorKind = C.CXCursor_UnionDecl
	/** \brief A C++ class. */
	Cursor_ClassDecl CursorKind = C.CXCursor_ClassDecl
	/** \brief An enumeration. */
	Cursor_EnumDecl CursorKind = C.CXCursor_EnumDecl
	/**
	 * \brief A field (in C) or non-static data member (in C++) in a
	 * struct, union, or C++ class.
	 */
	Cursor_FieldDecl CursorKind = C.CXCursor_FieldDecl
	/** \brief An enumerator constant. */
	Cursor_EnumConstantDecl CursorKind = C.CXCursor_EnumConstantDecl
	/** \brief A function. */
	Cursor_FunctionDecl CursorKind = C.CXCursor_FunctionDecl
	/** \brief A variable. */
	Cursor_VarDecl CursorKind = C.CXCursor_VarDecl
	/** \brief A function or method parameter. */
	Cursor_ParmDecl CursorKind = C.CXCursor_ParmDecl
	/** \brief An Objective-C \@interface. */
	Cursor_ObjCInterfaceDecl CursorKind = C.CXCursor_ObjCInterfaceDecl
	/** \brief An Objective-C \@interface for a category. */
	Cursor_ObjCCategoryDecl CursorKind = C.CXCursor_ObjCCategoryDecl
	/** \brief An Objective-C \@protocol declaration. */
	Cursor_ObjCProtocolDecl CursorKind = C.CXCursor_ObjCProtocolDecl
	/** \brief An Objective-C \@property declaration. */
	Cursor_ObjCPropertyDecl CursorKind = C.CXCursor_ObjCPropertyDecl
	/** \brief An Objective-C instance variable. */
	Cursor_ObjCIvarDecl CursorKind = C.CXCursor_ObjCIvarDecl
	/** \brief An Objective-C instance method. */
	Cursor_ObjCInstanceMethodDecl CursorKind = C.CXCursor_ObjCInstanceMethodDecl
	/** \brief An Objective-C class method. */
	Cursor_ObjCClassMethodDecl CursorKind = C.CXCursor_ObjCClassMethodDecl
	/** \brief An Objective-C \@implementation. */
	Cursor_ObjCImplementationDecl CursorKind = C.CXCursor_ObjCImplementationDecl
	/** \brief An Objective-C \@implementation for a category. */
	Cursor_ObjCCategoryImplDecl CursorKind = C.CXCursor_ObjCCategoryImplDecl
	/** \brief A typedef */
	Cursor_TypedefDecl CursorKind = C.CXCursor_TypedefDecl
	/** \brief A C++ class method. */
	Cursor_CXXMethod CursorKind = C.CXCursor_CXXMethod
	/** \brief A C++ namespace. */
	Cursor_Namespace CursorKind = C.CXCursor_Namespace
	/** \brief A linkage specification, e.g. 'extern "C"'. */
	Cursor_LinkageSpec CursorKind = C.CXCursor_LinkageSpec
	/** \brief A C++ constructor. */
	Cursor_Constructor CursorKind = C.CXCursor_Constructor
	/** \brief A C++ destructor. */
	Cursor_Destructor CursorKind = C.CXCursor_Destructor
	/** \brief A C++ conversion function. */
	Cursor_ConversionFunction CursorKind = C.CXCursor_ConversionFunction
	/** \brief A C++ template type parameter. */
	Cursor_TemplateTypeParameter CursorKind = C.CXCursor_TemplateTypeParameter
	/** \brief A C++ non-type template parameter. */
	Cursor_NonTypeTemplateParameter CursorKind = C.CXCursor_NonTypeTemplateParameter
	/** \brief A C++ template template parameter. */
	Cursor_TemplateTemplateParameter CursorKind = C.CXCursor_TemplateTemplateParameter
	/** \brief A C++ function template. */
	Cursor_FunctionTemplate CursorKind = C.CXCursor_FunctionTemplate
	/** \brief A C++ class template. */
	Cursor_ClassTemplate CursorKind = C.CXCursor_ClassTemplate
	/** \brief A C++ class template partial specialization. */
	Cursor_ClassTemplatePartialSpecialization CursorKind = C.CXCursor_ClassTemplatePartialSpecialization
	/** \brief A C++ namespace alias declaration. */
	Cursor_NamespaceAlias CursorKind = C.CXCursor_NamespaceAlias
	/** \brief A C++ using directive. */
	Cursor_UsingDirective CursorKind = C.CXCursor_UsingDirective
	/** \brief A C++ using declaration. */
	Cursor_UsingDeclaration CursorKind = C.CXCursor_UsingDeclaration
	/** \brief A C++ alias declaration */
	Cursor_TypeAliasDecl CursorKind = C.CXCursor_TypeAliasDecl
	/** \brief An Objective-C \@synthesize definition. */
	Cursor_ObjCSynthesizeDecl CursorKind = C.CXCursor_ObjCSynthesizeDecl
	/** \brief An Objective-C \@dynamic definition. */
	Cursor_ObjCDynamicDecl CursorKind = C.CXCursor_ObjCDynamicDecl
	/** \brief An access specifier. */
	Cursor_CXXAccessSpecifier CursorKind = C.CXCursor_CXXAccessSpecifier
	/** \brief An access specifier. */
	Cursor_FirstDecl CursorKind = C.CXCursor_FirstDecl
	/** \brief An access specifier. */
	Cursor_LastDecl CursorKind = C.CXCursor_LastDecl
	/** \brief An access specifier. */
	Cursor_FirstRef CursorKind = C.CXCursor_FirstRef
	/** \brief An access specifier. */
	Cursor_ObjCSuperClassRef CursorKind = C.CXCursor_ObjCSuperClassRef
	/** \brief An access specifier. */
	Cursor_ObjCProtocolRef CursorKind = C.CXCursor_ObjCProtocolRef
	/** \brief An access specifier. */
	Cursor_ObjCClassRef CursorKind = C.CXCursor_ObjCClassRef
	/**
	 * \brief A reference to a type declaration.
	 *
	 * A type reference occurs anywhere where a type is named but not
	 * declared. For example, given:
	 *
	 * \code
	 * typedef unsigned size_type;
	 * size_type size;
	 * \endcode
	 *
	 * The typedef is a declaration of size_type (CXCursor_TypedefDecl),
	 * while the type of the variable "size" is referenced. The cursor
	 * referenced by the type of size is the typedef for size_type.
	 */
	Cursor_TypeRef CursorKind = C.CXCursor_TypeRef
	/**
	 * \brief A reference to a type declaration.
	 *
	 * A type reference occurs anywhere where a type is named but not
	 * declared. For example, given:
	 *
	 * \code
	 * typedef unsigned size_type;
	 * size_type size;
	 * \endcode
	 *
	 * The typedef is a declaration of size_type (CXCursor_TypedefDecl),
	 * while the type of the variable "size" is referenced. The cursor
	 * referenced by the type of size is the typedef for size_type.
	 */
	Cursor_CXXBaseSpecifier CursorKind = C.CXCursor_CXXBaseSpecifier
	/**
	 * \brief A reference to a class template, function template, template
	 * template parameter, or class template partial specialization.
	 */
	Cursor_TemplateRef CursorKind = C.CXCursor_TemplateRef
	/**
	 * \brief A reference to a namespace or namespace alias.
	 */
	Cursor_NamespaceRef CursorKind = C.CXCursor_NamespaceRef
	/**
	 * \brief A reference to a member of a struct, union, or class that occurs in
	 * some non-expression context, e.g., a designated initializer.
	 */
	Cursor_MemberRef CursorKind = C.CXCursor_MemberRef
	/**
	 * \brief A reference to a labeled statement.
	 *
	 * This cursor kind is used to describe the jump to "start_over" in the
	 * goto statement in the following example:
	 *
	 * \code
	 *   start_over:
	 *     ++counter;
	 *
	 *     goto start_over;
	 * \endcode
	 *
	 * A label reference cursor refers to a label statement.
	 */
	Cursor_LabelRef CursorKind = C.CXCursor_LabelRef
	/**
	 * \brief A reference to a set of overloaded functions or function templates
	 * that has not yet been resolved to a specific function or function template.
	 *
	 * An overloaded declaration reference cursor occurs in C++ templates where
	 * a dependent name refers to a function. For example:
	 *
	 * \code
	 * template<typename T> void swap(T&, T&);
	 *
	 * struct X { ... };
	 * void swap(X&, X&);
	 *
	 * template<typename T>
	 * void reverse(T* first, T* last) {
	 *   while (first < last - 1) {
	 *     swap(*first, *--last);
	 *     ++first;
	 *   }
	 * }
	 *
	 * struct Y { };
	 * void swap(Y&, Y&);
	 * \endcode
	 *
	 * Here, the identifier "swap" is associated with an overloaded declaration
	 * reference. In the template definition, "swap" refers to either of the two
	 * "swap" functions declared above, so both results will be available. At
	 * instantiation time, "swap" may also refer to other functions found via
	 * argument-dependent lookup (e.g., the "swap" function at the end of the
	 * example).
	 *
	 * The functions \c clang_getNumOverloadedDecls() and
	 * \c clang_getOverloadedDecl() can be used to retrieve the definitions
	 * referenced by this cursor.
	 */
	Cursor_OverloadedDeclRef CursorKind = C.CXCursor_OverloadedDeclRef
	/**
	 * \brief A reference to a variable that occurs in some non-expression
	 * context, e.g., a C++ lambda capture list.
	 */
	Cursor_VariableRef CursorKind = C.CXCursor_VariableRef
	/**
	 * \brief A reference to a variable that occurs in some non-expression
	 * context, e.g., a C++ lambda capture list.
	 */
	Cursor_LastRef CursorKind = C.CXCursor_LastRef
	/**
	 * \brief A reference to a variable that occurs in some non-expression
	 * context, e.g., a C++ lambda capture list.
	 */
	Cursor_FirstInvalid CursorKind = C.CXCursor_FirstInvalid
	/**
	 * \brief A reference to a variable that occurs in some non-expression
	 * context, e.g., a C++ lambda capture list.
	 */
	Cursor_InvalidFile CursorKind = C.CXCursor_InvalidFile
	/**
	 * \brief A reference to a variable that occurs in some non-expression
	 * context, e.g., a C++ lambda capture list.
	 */
	Cursor_NoDeclFound CursorKind = C.CXCursor_NoDeclFound
	/**
	 * \brief A reference to a variable that occurs in some non-expression
	 * context, e.g., a C++ lambda capture list.
	 */
	Cursor_NotImplemented CursorKind = C.CXCursor_NotImplemented
	/**
	 * \brief A reference to a variable that occurs in some non-expression
	 * context, e.g., a C++ lambda capture list.
	 */
	Cursor_InvalidCode CursorKind = C.CXCursor_InvalidCode
	/**
	 * \brief A reference to a variable that occurs in some non-expression
	 * context, e.g., a C++ lambda capture list.
	 */
	Cursor_LastInvalid CursorKind = C.CXCursor_LastInvalid
	/**
	 * \brief A reference to a variable that occurs in some non-expression
	 * context, e.g., a C++ lambda capture list.
	 */
	Cursor_FirstExpr CursorKind = C.CXCursor_FirstExpr
	/**
	 * \brief An expression whose specific kind is not exposed via this
	 * interface.
	 *
	 * Unexposed expressions have the same operations as any other kind
	 * of expression; one can extract their location information,
	 * spelling, children, etc. However, the specific kind of the
	 * expression is not reported.
	 */
	Cursor_UnexposedExpr CursorKind = C.CXCursor_UnexposedExpr
	/**
	 * \brief An expression that refers to some value declaration, such
	 * as a function, varible, or enumerator.
	 */
	Cursor_DeclRefExpr CursorKind = C.CXCursor_DeclRefExpr
	/**
	 * \brief An expression that refers to a member of a struct, union,
	 * class, Objective-C class, etc.
	 */
	Cursor_MemberRefExpr CursorKind = C.CXCursor_MemberRefExpr
	/** \brief An expression that calls a function. */
	Cursor_CallExpr CursorKind = C.CXCursor_CallExpr
	/** \brief An expression that sends a message to an Objective-C
	  object or class. */
	Cursor_ObjCMessageExpr CursorKind = C.CXCursor_ObjCMessageExpr
	/** \brief An expression that represents a block literal. */
	Cursor_BlockExpr CursorKind = C.CXCursor_BlockExpr
	/** \brief An integer literal.
	 */
	Cursor_IntegerLiteral CursorKind = C.CXCursor_IntegerLiteral
	/** \brief A floating point number literal.
	 */
	Cursor_FloatingLiteral CursorKind = C.CXCursor_FloatingLiteral
	/** \brief An imaginary number literal.
	 */
	Cursor_ImaginaryLiteral CursorKind = C.CXCursor_ImaginaryLiteral
	/** \brief A string literal.
	 */
	Cursor_StringLiteral CursorKind = C.CXCursor_StringLiteral
	/** \brief A character literal.
	 */
	Cursor_CharacterLiteral CursorKind = C.CXCursor_CharacterLiteral
	/** \brief A parenthesized expression, e.g. "(1)".
	 *
	 * This AST node is only formed if full location information is requested.
	 */
	Cursor_ParenExpr CursorKind = C.CXCursor_ParenExpr
	/** \brief This represents the unary-expression's (except sizeof and
	 * alignof).
	 */
	Cursor_UnaryOperator CursorKind = C.CXCursor_UnaryOperator
	/** \brief [C99 6.5.2.1] Array Subscripting.
	 */
	Cursor_ArraySubscriptExpr CursorKind = C.CXCursor_ArraySubscriptExpr
	/** \brief A builtin binary operation expression such as "x + y" or
	 * "x <= y".
	 */
	Cursor_BinaryOperator CursorKind = C.CXCursor_BinaryOperator
	/** \brief Compound assignment such as "+=".
	 */
	Cursor_CompoundAssignOperator CursorKind = C.CXCursor_CompoundAssignOperator
	/** \brief The ?: ternary operator.
	 */
	Cursor_ConditionalOperator CursorKind = C.CXCursor_ConditionalOperator
	/** \brief An explicit cast in C (C99 6.5.4) or a C-style cast in C++
	 * (C++ [expr.cast]), which uses the syntax (Type)expr.
	 *
	 * For example: (int)f.
	 */
	Cursor_CStyleCastExpr CursorKind = C.CXCursor_CStyleCastExpr
	/** \brief [C99 6.5.2.5]
	 */
	Cursor_CompoundLiteralExpr CursorKind = C.CXCursor_CompoundLiteralExpr
	/** \brief Describes an C or C++ initializer list.
	 */
	Cursor_InitListExpr CursorKind = C.CXCursor_InitListExpr
	/** \brief The GNU address of label extension, representing &&label.
	 */
	Cursor_AddrLabelExpr CursorKind = C.CXCursor_AddrLabelExpr
	/** \brief This is the GNU Statement Expression extension: ({int X=4; X;})
	 */
	Cursor_StmtExpr CursorKind = C.CXCursor_StmtExpr
	/** \brief Represents a C11 generic selection.
	 */
	Cursor_GenericSelectionExpr CursorKind = C.CXCursor_GenericSelectionExpr
	/** \brief Implements the GNU __null extension, which is a name for a null
	 * pointer constant that has integral type (e.g., int or long) and is the same
	 * size and alignment as a pointer.
	 *
	 * The __null extension is typically only used by system headers, which define
	 * NULL as __null in C++ rather than using 0 (which is an integer that may not
	 * match the size of a pointer).
	 */
	Cursor_GNUNullExpr CursorKind = C.CXCursor_GNUNullExpr
	/** \brief C++'s static_cast<> expression.
	 */
	Cursor_CXXStaticCastExpr CursorKind = C.CXCursor_CXXStaticCastExpr
	/** \brief C++'s dynamic_cast<> expression.
	 */
	Cursor_CXXDynamicCastExpr CursorKind = C.CXCursor_CXXDynamicCastExpr
	/** \brief C++'s reinterpret_cast<> expression.
	 */
	Cursor_CXXReinterpretCastExpr CursorKind = C.CXCursor_CXXReinterpretCastExpr
	/** \brief C++'s const_cast<> expression.
	 */
	Cursor_CXXConstCastExpr CursorKind = C.CXCursor_CXXConstCastExpr
	/** \brief Represents an explicit C++ type conversion that uses "functional"
	 * notion (C++ [expr.type.conv]).
	 *
	 * Example:
	 * \code
	 *   x = int(0.5);
	 * \endcode
	 */
	Cursor_CXXFunctionalCastExpr CursorKind = C.CXCursor_CXXFunctionalCastExpr
	/** \brief A C++ typeid expression (C++ [expr.typeid]).
	 */
	Cursor_CXXTypeidExpr CursorKind = C.CXCursor_CXXTypeidExpr
	/** \brief [C++ 2.13.5] C++ Boolean Literal.
	 */
	Cursor_CXXBoolLiteralExpr CursorKind = C.CXCursor_CXXBoolLiteralExpr
	/** \brief [C++0x 2.14.7] C++ Pointer Literal.
	 */
	Cursor_CXXNullPtrLiteralExpr CursorKind = C.CXCursor_CXXNullPtrLiteralExpr
	/** \brief Represents the "this" expression in C++
	 */
	Cursor_CXXThisExpr CursorKind = C.CXCursor_CXXThisExpr
	/** \brief [C++ 15] C++ Throw Expression.
	 *
	 * This handles 'throw' and 'throw' assignment-expression. When
	 * assignment-expression isn't present, Op will be null.
	 */
	Cursor_CXXThrowExpr CursorKind = C.CXCursor_CXXThrowExpr
	/** \brief A new expression for memory allocation and constructor calls, e.g:
	 * "new CXXNewExpr(foo)".
	 */
	Cursor_CXXNewExpr CursorKind = C.CXCursor_CXXNewExpr
	/** \brief A delete expression for memory deallocation and destructor calls,
	 * e.g. "delete[] pArray".
	 */
	Cursor_CXXDeleteExpr CursorKind = C.CXCursor_CXXDeleteExpr
	/** \brief A unary expression.
	 */
	Cursor_UnaryExpr CursorKind = C.CXCursor_UnaryExpr
	/** \brief An Objective-C string literal i.e. @"foo".
	 */
	Cursor_ObjCStringLiteral CursorKind = C.CXCursor_ObjCStringLiteral
	/** \brief An Objective-C \@encode expression.
	 */
	Cursor_ObjCEncodeExpr CursorKind = C.CXCursor_ObjCEncodeExpr
	/** \brief An Objective-C \@selector expression.
	 */
	Cursor_ObjCSelectorExpr CursorKind = C.CXCursor_ObjCSelectorExpr
	/** \brief An Objective-C \@protocol expression.
	 */
	Cursor_ObjCProtocolExpr CursorKind = C.CXCursor_ObjCProtocolExpr
	/** \brief An Objective-C "bridged" cast expression, which casts between
	 * Objective-C pointers and C pointers, transferring ownership in the process.
	 *
	 * \code
	 *   NSString *str = (__bridge_transfer NSString *)CFCreateString();
	 * \endcode
	 */
	Cursor_ObjCBridgedCastExpr CursorKind = C.CXCursor_ObjCBridgedCastExpr
	/** \brief Represents a C++0x pack expansion that produces a sequence of
	 * expressions.
	 *
	 * A pack expansion expression contains a pattern (which itself is an
	 * expression) followed by an ellipsis. For example:
	 *
	 * \code
	 * template<typename F, typename ...Types>
	 * void forward(F f, Types &&...args) {
	 *  f(static_cast<Types&&>(args)...);
	 * }
	 * \endcode
	 */
	Cursor_PackExpansionExpr CursorKind = C.CXCursor_PackExpansionExpr
	/** \brief Represents an expression that computes the length of a parameter
	 * pack.
	 *
	 * \code
	 * template<typename ...Types>
	 * struct count {
	 *   static const unsigned value = sizeof...(Types);
	 * };
	 * \endcode
	 */
	Cursor_SizeOfPackExpr CursorKind = C.CXCursor_SizeOfPackExpr
	Cursor_LambdaExpr     CursorKind = C.CXCursor_LambdaExpr
	/** \brief Objective-c Boolean Literal.
	 */
	Cursor_ObjCBoolLiteralExpr CursorKind = C.CXCursor_ObjCBoolLiteralExpr
	/** \brief Represents the "self" expression in a ObjC method.
	 */
	Cursor_ObjCSelfExpr CursorKind = C.CXCursor_ObjCSelfExpr
	/** \brief Represents the "self" expression in a ObjC method.
	 */
	Cursor_LastExpr CursorKind = C.CXCursor_LastExpr
	/** \brief Represents the "self" expression in a ObjC method.
	 */
	Cursor_FirstStmt CursorKind = C.CXCursor_FirstStmt
	/**
	 * \brief A statement whose specific kind is not exposed via this
	 * interface.
	 *
	 * Unexposed statements have the same operations as any other kind of
	 * statement; one can extract their location information, spelling,
	 * children, etc. However, the specific kind of the statement is not
	 * reported.
	 */
	Cursor_UnexposedStmt CursorKind = C.CXCursor_UnexposedStmt
	/** \brief A labelled statement in a function.
	 *
	 * This cursor kind is used to describe the "start_over:" label statement in
	 * the following example:
	 *
	 * \code
	 *   start_over:
	 *     ++counter;
	 * \endcode
	 *
	 */
	Cursor_LabelStmt CursorKind = C.CXCursor_LabelStmt
	/** \brief A group of statements like { stmt stmt }.
	 *
	 * This cursor kind is used to describe compound statements, e.g. function
	 * bodies.
	 */
	Cursor_CompoundStmt CursorKind = C.CXCursor_CompoundStmt
	/** \brief A case statement.
	 */
	Cursor_CaseStmt CursorKind = C.CXCursor_CaseStmt
	/** \brief A default statement.
	 */
	Cursor_DefaultStmt CursorKind = C.CXCursor_DefaultStmt
	/** \brief An if statement
	 */
	Cursor_IfStmt CursorKind = C.CXCursor_IfStmt
	/** \brief A switch statement.
	 */
	Cursor_SwitchStmt CursorKind = C.CXCursor_SwitchStmt
	/** \brief A while statement.
	 */
	Cursor_WhileStmt CursorKind = C.CXCursor_WhileStmt
	/** \brief A do statement.
	 */
	Cursor_DoStmt CursorKind = C.CXCursor_DoStmt
	/** \brief A for statement.
	 */
	Cursor_ForStmt CursorKind = C.CXCursor_ForStmt
	/** \brief A goto statement.
	 */
	Cursor_GotoStmt CursorKind = C.CXCursor_GotoStmt
	/** \brief An indirect goto statement.
	 */
	Cursor_IndirectGotoStmt CursorKind = C.CXCursor_IndirectGotoStmt
	/** \brief A continue statement.
	 */
	Cursor_ContinueStmt CursorKind = C.CXCursor_ContinueStmt
	/** \brief A break statement.
	 */
	Cursor_BreakStmt CursorKind = C.CXCursor_BreakStmt
	/** \brief A return statement.
	 */
	Cursor_ReturnStmt CursorKind = C.CXCursor_ReturnStmt
	/** \brief A GCC inline assembly statement extension.
	 */
	Cursor_GCCAsmStmt CursorKind = C.CXCursor_GCCAsmStmt
	/** \brief A GCC inline assembly statement extension.
	 */
	Cursor_AsmStmt CursorKind = C.CXCursor_AsmStmt
	/** \brief Objective-C's overall \@try-\@catch-\@finally statement.
	 */
	Cursor_ObjCAtTryStmt CursorKind = C.CXCursor_ObjCAtTryStmt
	/** \brief Objective-C's \@catch statement.
	 */
	Cursor_ObjCAtCatchStmt CursorKind = C.CXCursor_ObjCAtCatchStmt
	/** \brief Objective-C's \@finally statement.
	 */
	Cursor_ObjCAtFinallyStmt CursorKind = C.CXCursor_ObjCAtFinallyStmt
	/** \brief Objective-C's \@throw statement.
	 */
	Cursor_ObjCAtThrowStmt CursorKind = C.CXCursor_ObjCAtThrowStmt
	/** \brief Objective-C's \@synchronized statement.
	 */
	Cursor_ObjCAtSynchronizedStmt CursorKind = C.CXCursor_ObjCAtSynchronizedStmt
	/** \brief Objective-C's autorelease pool statement.
	 */
	Cursor_ObjCAutoreleasePoolStmt CursorKind = C.CXCursor_ObjCAutoreleasePoolStmt
	/** \brief Objective-C's collection statement.
	 */
	Cursor_ObjCForCollectionStmt CursorKind = C.CXCursor_ObjCForCollectionStmt
	/** \brief C++'s catch statement.
	 */
	Cursor_CXXCatchStmt CursorKind = C.CXCursor_CXXCatchStmt
	/** \brief C++'s try statement.
	 */
	Cursor_CXXTryStmt CursorKind = C.CXCursor_CXXTryStmt
	/** \brief C++'s for (* : *) statement.
	 */
	Cursor_CXXForRangeStmt CursorKind = C.CXCursor_CXXForRangeStmt
	/** \brief Windows Structured Exception Handling's try statement.
	 */
	Cursor_SEHTryStmt CursorKind = C.CXCursor_SEHTryStmt
	/** \brief Windows Structured Exception Handling's except statement.
	 */
	Cursor_SEHExceptStmt CursorKind = C.CXCursor_SEHExceptStmt
	/** \brief Windows Structured Exception Handling's finally statement.
	 */
	Cursor_SEHFinallyStmt CursorKind = C.CXCursor_SEHFinallyStmt
	/** \brief A MS inline assembly statement extension.
	 */
	Cursor_MSAsmStmt CursorKind = C.CXCursor_MSAsmStmt
	/** \brief The null satement ";": C99 6.8.3p3.
	 *
	 * This cursor kind is used to describe the null statement.
	 */
	Cursor_NullStmt CursorKind = C.CXCursor_NullStmt
	/** \brief Adaptor class for mixing declarations with statements and
	 * expressions.
	 */
	Cursor_DeclStmt CursorKind = C.CXCursor_DeclStmt
	/** \brief OpenMP parallel directive.
	 */
	Cursor_OMPParallelDirective CursorKind = C.CXCursor_OMPParallelDirective
	/** \brief OpenMP parallel directive.
	 */
	Cursor_LastStmt CursorKind = C.CXCursor_LastStmt
	/**
	 * \brief Cursor that represents the translation unit itself.
	 *
	 * The translation unit cursor exists primarily to act as the root
	 * cursor for traversing the contents of a translation unit.
	 */
	Cursor_TranslationUnit CursorKind = C.CXCursor_TranslationUnit
	/**
	 * \brief Cursor that represents the translation unit itself.
	 *
	 * The translation unit cursor exists primarily to act as the root
	 * cursor for traversing the contents of a translation unit.
	 */
	Cursor_FirstAttr CursorKind = C.CXCursor_FirstAttr
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_UnexposedAttr CursorKind = C.CXCursor_UnexposedAttr
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_IBActionAttr CursorKind = C.CXCursor_IBActionAttr
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_IBOutletAttr CursorKind = C.CXCursor_IBOutletAttr
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_IBOutletCollectionAttr CursorKind = C.CXCursor_IBOutletCollectionAttr
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_CXXFinalAttr CursorKind = C.CXCursor_CXXFinalAttr
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_CXXOverrideAttr CursorKind = C.CXCursor_CXXOverrideAttr
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_AnnotateAttr CursorKind = C.CXCursor_AnnotateAttr
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_AsmLabelAttr CursorKind = C.CXCursor_AsmLabelAttr
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_PackedAttr CursorKind = C.CXCursor_PackedAttr
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_LastAttr CursorKind = C.CXCursor_LastAttr
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_PreprocessingDirective CursorKind = C.CXCursor_PreprocessingDirective
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_MacroDefinition CursorKind = C.CXCursor_MacroDefinition
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_MacroExpansion CursorKind = C.CXCursor_MacroExpansion
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_MacroInstantiation CursorKind = C.CXCursor_MacroInstantiation
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_InclusionDirective CursorKind = C.CXCursor_InclusionDirective
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_FirstPreprocessing CursorKind = C.CXCursor_FirstPreprocessing
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	Cursor_LastPreprocessing CursorKind = C.CXCursor_LastPreprocessing
	/**
	 * \brief A module import declaration.
	 */
	Cursor_ModuleImportDecl CursorKind = C.CXCursor_ModuleImportDecl
	/**
	 * \brief A module import declaration.
	 */
	Cursor_FirstExtraDecl CursorKind = C.CXCursor_FirstExtraDecl
	/**
	 * \brief A module import declaration.
	 */
	Cursor_LastExtraDecl CursorKind = C.CXCursor_LastExtraDecl
)
