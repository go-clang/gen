package gen

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
)

const fakeStatement = "REMOVE"

// ASTFunc represents a AST Func.
type ASTFunc struct {
	*ast.FuncDecl

	f   *Function
	ret *ast.ReturnStmt
}

// NewASTFunc returns the new initialized ASTFunc.
func NewASTFunc(f *Function) *ASTFunc {
	return &ASTFunc{
		FuncDecl: &ast.FuncDecl{
			Name: &ast.Ident{
				Name: f.Name,
			},
			Type: &ast.FuncType{
				Results: &ast.FieldList{
					List: []*ast.Field{},
				},
			},
			Body: &ast.BlockStmt{},
		},
		f: f,
		ret: &ast.ReturnStmt{
			Results: []ast.Expr{},
		},
	}
}

// GenerateFunctionString generates function string.
func GenerateFunctionString(af *ASTFunc) string {
	var b strings.Builder
	if err := format.Node(&b, token.NewFileSet(), []ast.Decl{af.FuncDecl}); err != nil {
		panic(fmt.Errorf("unexpected error: %w", err))
	}

	fnName := b.String()
	fnName = strings.ReplaceAll(fnName, fmt.Sprintf("%s()", fakeStatement), "")

	return fnName
}

// Generate generates function.
func (af *ASTFunc) Generate() {
	// TODO(go-clang): maybe name the return arguments...
	// because of clang_getDiagnosticOption -> the normal return can be always just "o"?
	// https://github.com/go-clang/gen/issues/57

	// TODO(go-clang): re-enable this, see the comment at the bottom of the generate function
	// for details https://github.com/go-clang/gen/issues/54
	// Add function comment
	//
	// if f.Comment != "" {
	// 	fa.Doc = &ast.CommentGroup{
	// 		List: []*ast.Comment{
	// 			&ast.Comment{
	// 				Text: f.Comment,
	// 			},
	// 		},
	// 	}
	// }
	af.GenerateReceiver()

	if af.f.Member != nil {
		if af.f.ReturnType.IsSlice {
			af.AddStatement(doDeclare("s", doGoType(af.f.ReturnType)))
			af.AddCToGoSliceConversion("s", af.f.Receiver.Name+".c."+af.f.Member.Name, af.f.Receiver.Name+".c."+af.f.ReturnType.LengthOfSlice)

			af.AddReturnItem(&ast.Ident{
				Name: "s",
			})
			af.AddReturnType("", af.f.ReturnType)
			af.AddEmptyLine()

			// add the return statement
			af.AddStatement(af.ret)
		} else {
			af.GenerateReturn(&ast.SelectorExpr{
				X: accessMember(af.f.Receiver.Name, "c"),
				Sel: &ast.Ident{
					Name: af.f.Member.Name,
				},
			})
		}
	} else {
		// basic call to the C function
		call := doCCast(af.f.CName)

		if callArguments := af.GenerateParameters(); len(callArguments) > 0 {
			call.Args = callArguments
		}

		af.GenerateReturn(call)
	}
}

// GenerateReceiver generates function receiver.
func (af *ASTFunc) GenerateReceiver() {
	// add receiver to make function a method
	if af.f.Receiver.Name == "" {
		return
	}

	// TODO(go-clang): maybe to not set the receiver at all? -> do this outside of the generate function?
	// https://github.com/go-clang/gen/issues/52
	if len(af.f.Parameters) > 0 {
		af.Recv = &ast.FieldList{
			List: []*ast.Field{
				doField(af.f.Receiver.Name, af.f.Receiver.Type),
			},
		}
	}
}

// GenerateParameters generates function parameters.
func (af *ASTFunc) GenerateParameters() []ast.Expr {
	if len(af.f.Parameters) == 0 {
		return nil
	}

	var callArguments []ast.Expr

	if af.f.Receiver.Name != "" {
		af.f.Parameters[0].Name = af.f.Receiver.Name
	}

	af.Type.Params = &ast.FieldList{
		List: []*ast.Field{},
	}

	hasDeclaration := false

	// add parameters to the function
	for i, p := range af.f.Parameters {
		if i == 0 && af.f.Receiver.Name != "" {
			continue
		}

		// ignore length parameters since they will be filled by the slice itself
		if p.Type.LengthOfSlice != "" && !p.Type.IsReturnArgument {
			continue
		}

		switch {
		case p.Type.IsSlice && !p.Type.IsReturnArgument:
			hasDeclaration = true

			if p.Type.CGoName == CSChar && p.Type.PointerLevel >= 1 { // one pointer level from being a string, one from being an array
				af.AddGoToCSliceConversion(p.Name, p.Type)
			} else {
				af.AddCArrayFromGoSlice(p.Name, p.Type)
			}

		case p.Type.IsReturnArgument:
			if p.Type.LengthOfSlice == "" {
				// add the return type to the function return arguments
				retType := p.Type
				if p.Type.GoName == "cxstring" {
					retType.GoName = "string"
				}

				af.AddReturnType("", retType)
			}

			if p.Type.IsSlice && p.Type.IsReturnArgument {
				af.AddStatement(doDeclare(
					"cp_"+p.Name,
					sliceType(p.Type),
				))
			}

			// declare the return argument's variable
			var varType ast.Expr
			switch {
			case p.Type.PointerLevel > 0 && p.Type.CGoName == CSChar:
				varType = doPointer(doCType("char"))

			case p.Type.IsPrimitive:
				varType = doCType(p.Type.CGoName)

			default:
				varType = &ast.Ident{
					Name: p.Type.GoName,
				}
			}
			if p.Type.IsSlice {
				varType = &ast.ArrayType{
					Elt: varType,
				}
			}

			af.AddStatement(doDeclare(
				p.Name,
				varType,
			))

			switch {
			case p.Type.GoName == "cxstring":
				af.AddDefer(doCall(p.Name, "Dispose"))

			case p.Type.PointerLevel > 0 && p.Type.CGoName == CSChar:
				af.AddDefer(doCCast(
					"free",
					doCall(
						"unsafe",
						"Pointer",
						&ast.Ident{
							Name: p.Name,
						},
					),
				))
			}

			if p.Type.LengthOfSlice == "" {
				// add the return argument to the return statement
				switch {
				case p.Type.PointerLevel > 0 && p.Type.CGoName == CSChar:
					af.AddReturnItem(doCCast(
						"GoString",
						&ast.Ident{
							Name: p.Name,
						},
					))

				case p.Type.GoName == "cxstring":
					af.AddReturnItem(doCall(p.Name, "String"))

				case p.Type.IsPrimitive:
					af.AddReturnItem(doCast(
						p.Type.GoName,
						&ast.Ident{
							Name: p.Name,
						},
					))
				default:
					af.AddReturnItem(&ast.Ident{
						Name: p.Name,
					})
				}
			}
			continue

		case p.Type.PointerLevel > 0 && p.Type.IsPrimitive && p.Type.CGoName != CSChar:
			hasDeclaration = true

			varType := doCType(p.Type.CGoName)

			af.AddStatement(doDeclare(
				"cp_"+p.Name,
				varType,
			))
			af.AddStatement(&ast.IfStmt{
				Cond: &ast.BinaryExpr{
					X: &ast.Ident{
						Name: p.Name,
					},
					Op: token.NEQ,
					Y: &ast.Ident{
						Name: "nil",
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "cp_" + p.Name,
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: varType,
									Args: []ast.Expr{
										doPointer(&ast.Ident{
											Name: p.Name,
										}),
									},
								},
							},
						},
					},
				},
			})
		}

		af.Type.Params.List = append(af.Type.Params.List, doField(p.Name, p.Type))
	}

	if len(af.ret.Results) > 0 || hasDeclaration {
		af.AddEmptyLine()
	}

	goToCTypeConversions := false

	// add arguments to the C function call
	for _, p := range af.f.Parameters {
		var pf ast.Expr

		switch {
		case p.Type.IsSlice:
			pf = &ast.Ident{
				Name: "cp_" + p.Name,
			}

		case p.Type.IsPrimitive:
			// handle Go type to C type conversions
			switch {
			case p.Type.PointerLevel == 1 && p.Type.CGoName == CSChar:
				goToCTypeConversions = true

				af.AddAssignment(
					"c_"+p.Name,
					doCCast(
						"CString",
						&ast.Ident{
							Name: p.Name,
						},
					),
				)
				af.AddDefer(doCCast(
					"free",
					doCall(
						"unsafe",
						"Pointer",
						&ast.Ident{
							Name: "c_" + p.Name,
						},
					),
				))

				pf = &ast.Ident{
					Name: "c_" + p.Name,
				}

			// TODO(go-clang): try to get cxstring and "String" completely out of this function since it is just a struct which can be handled by the struct code
			// https://github.com/go-clang/gen/issues/25
			case p.Type.CGoName == "cxstring":
				pf = accessMember(p.Name, "c")

			default:
				switch {
				case p.Type.IsReturnArgument:
					// return arguments already have a cast
					pf = &ast.Ident{
						Name: p.Name,
					}

				case p.Type.LengthOfSlice != "":
					pf = doCCast(
						p.Type.CGoName,
						doCast(
							"len",
							&ast.Ident{
								Name: p.Type.LengthOfSlice,
							},
						),
					)

				case p.Type.PointerLevel > 0:
					pf = doReference(&ast.Ident{
						Name: "cp_" + p.Name,
					})

				default:
					pf = doCCast(
						p.Type.CGoName,
						&ast.Ident{
							Name: p.Name,
						},
					)
				}
			}

		default:
			pf = accessMember(p.Name, "c")

			if p.Type.PointerLevel > 0 && !p.Type.IsReturnArgument && !p.Type.IsPointerComposition {
				pf = doReference(pf)
			}
		}

		if p.Type.IsReturnArgument && !p.Type.IsPointerComposition {
			pf = doReference(pf)
		}

		callArguments = append(callArguments, pf)
	}

	if goToCTypeConversions {
		af.AddEmptyLine()
	}

	return callArguments
}

// GenerateReturn generates return statement.
func (af *ASTFunc) GenerateReturn(call ast.Expr) {
	returnType := af.f.ReturnType

	// check if we need to add a return
	if returnType.GoName != "void" || len(af.ret.Results) > 0 {
		if returnType.GoName == "cxstring" {
			// do the C function call and save the result into the new variable "o" while transforming it into a cxstring
			af.AddAssignment("o", doCompose("cxstring", call))
			af.AddDefer(doCall("o", "Dispose"))
			af.AddEmptyLine()

			// call the String method on the cxstring instance
			af.AddReturnItem(doCall("o", "String"))

			// change the return type to "string"
			af.AddReturnType("", Type{
				GoName: "string",
			})
		} else {
			if returnType.GoName != "void" {
				// add the function return type
				af.AddReturnType("", returnType)
			}

			// do we need to convert the return of the C function into a boolean?
			switch {
			case returnType.GoName == "bool":
				// do the C function call and save the result into the new variable "o"
				af.AddAssignment("o", call)
				af.AddEmptyLine()

				// check if o is not equal to zero and return the result
				af.AddReturnItem(&ast.BinaryExpr{
					X: &ast.Ident{
						Name: "o",
					},
					Op: token.NEQ,
					Y: doCCast(
						returnType.CGoName,
						doZero(),
					),
				})

			// TODO(go-clang): refactor the const char * check so that one function is used everywhere to check for that C type
			// https://github.com/go-clang/gen/issues/56
			case returnType.CGoName == CSChar && returnType.PointerLevel == 1:
				// if this is a normal const char * C type there is not so much to do
				af.AddReturnItem(doCCast(
					"GoString",
					call,
				))

			case returnType.GoName == "time.Time":
				af.AddReturnItem(doCall(
					"time",
					"Unix",
					doCast("int64", call),
					doZero(),
				))

			case returnType.GoName == "void":
				// handle the case where the C function has no return argument but parameters that are return arguments

				// do the C function call
				af.AddStatement(&ast.ExprStmt{
					X: call,
				})
				af.AddEmptyLine()

			case returnType.PointerLevel > 0:
				// do the C function call and save the result into the new variable "o"
				af.AddAssignment("o", call)
				af.AddEmptyLine()

				af.AddStatement(doDeclare(
					"gop_o",
					doGoType(returnType),
				))
				var compositionValue ast.Expr = &ast.Ident{
					Name: "o",
				}
				if returnType.IsPointerComposition && returnType.PointerLevel == 0 {
					compositionValue = doReference(compositionValue)
				} else if !returnType.IsPointerComposition && returnType.PointerLevel > 0 {
					compositionValue = doUnreference(compositionValue)
				}
				af.AddStatement(&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.Ident{
							Name: "o",
						},
						Op: token.NEQ,
						Y: &ast.Ident{
							Name: "nil",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "gop_o",
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									doReference(doCompose(
										returnType.GoName,
										compositionValue,
									)),
								},
							},
						},
					},
				})
				af.AddEmptyLine()

				af.AddReturnItem(&ast.Ident{
					Name: "gop_o",
				})

			default:
				var convCall ast.Expr

				// structs are literals, everything else is a cast
				if returnType.IsPrimitive {
					convCall = doCast(returnType.GoName, call)
				} else {
					convCall = doCompose(returnType.GoName, call)
				}

				if len(af.ret.Results) > 0 {
					// do the C function call and save the result into the new variable "o"
					af.AddAssignment("o", convCall)
					af.AddEmptyLine()

					// add the C function call result to the return statement
					af.AddReturnItem(&ast.Ident{
						Name: "o",
					})
				} else {
					af.AddReturnItem(convCall)
				}
			}
		}

		af.AddCToGoConversions()

		// add the return statement
		af.AddStatement(af.ret)
	} else {
		af.AddCToGoConversions()

		// no return needed, just add the C function call
		af.AddStatement(&ast.ExprStmt{
			X: call,
		})
	}
}

// AddStatement adds stmt ast.Stmt to af.
func (af *ASTFunc) AddStatement(stmt ast.Stmt) {
	af.Body.List = append(af.Body.List, stmt)
}

// AddAssignment adds assignment to af.
func (af *ASTFunc) AddAssignment(variable string, e ast.Expr) {
	af.AddStatement(&ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{
				Name: variable,
			},
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			e,
		},
	})
}

// AddDefer adds deferred statement to af.
func (af *ASTFunc) AddDefer(call *ast.CallExpr) {
	af.AddStatement(&ast.DeferStmt{
		Call: call,
	})
}

// AddEmptyLine adds empty line to af.
func (af *ASTFunc) AddEmptyLine() {
	// TODO(go-clang): this should be done using something else than a fakeStatement.
	// https://github.com/go-clang/gen/issues/53
	af.AddStatement(&ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.Ident{
				Name: fakeStatement,
			},
		},
	})
}

// AddReturnItem adds return item to af.
func (af *ASTFunc) AddReturnItem(item ast.Expr) {
	af.ret.Results = append(af.ret.Results, item)
}

// AddReturnType adds return type to af.
func (af *ASTFunc) AddReturnType(name string, typ Type) {
	af.Type.Results.List = append(af.Type.Results.List, doField(name, typ))
}

// AddCToGoConversions adds C to Go conversions to af.
func (af *ASTFunc) AddCToGoConversions() {
	cToGoTypeConversions := false

	for _, p := range af.f.Parameters {
		if p.Type.IsSlice && p.Type.IsReturnArgument {
			cToGoTypeConversions = true

			var lengthOfSlice string
			for _, pl := range af.f.Parameters {
				if pl.Type.LengthOfSlice == p.Name {
					lengthOfSlice = pl.Name

					break
				}
			}

			af.AddCToGoSliceConversion(p.Name, "cp_"+p.Name, lengthOfSlice)
		}
	}

	if cToGoTypeConversions {
		af.AddEmptyLine()
	}
}

// AddCToGoSliceConversion adds C to Go slice conversion to af.
func (af *ASTFunc) AddCToGoSliceConversion(name string, cname string, lengthOfSlice string) {
	af.AddAssignment("gos_"+name, &ast.CallExpr{
		Fun: &ast.ParenExpr{
			X: doPointer(accessMember("reflect", "SliceHeader")),
		},
		Args: []ast.Expr{
			doCall(
				"unsafe",
				"Pointer",
				doReference(&ast.Ident{
					Name: name,
				}),
			),
		},
	})

	af.AddStatement(&ast.AssignStmt{
		Lhs: []ast.Expr{
			accessMember("gos_"+name, "Cap"),
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			doCast(
				"int",
				&ast.Ident{
					Name: lengthOfSlice,
				},
			),
		},
	})

	af.AddStatement(&ast.AssignStmt{
		Lhs: []ast.Expr{
			accessMember("gos_"+name, "Len"),
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			doCast(
				"int",
				&ast.Ident{
					Name: lengthOfSlice,
				},
			),
		},
	})

	af.AddStatement(&ast.AssignStmt{
		Lhs: []ast.Expr{
			accessMember("gos_"+name, "Data"),
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			doCast(
				"uintptr",
				doCall(
					"unsafe",
					"Pointer",
					&ast.Ident{
						Name: cname,
					},
				),
			),
		},
	})
}

// AddCArrayFromGoSlice adds C array from Go slice to af.
func (af *ASTFunc) AddCArrayFromGoSlice(name string, typ Type) {
	sliceType := sliceType(typ)

	af.AddAssignment("gos_"+name, &ast.CallExpr{
		Fun: &ast.ParenExpr{
			X: doPointer(accessMember("reflect", "SliceHeader")),
		},
		Args: []ast.Expr{
			doCall(
				"unsafe",
				"Pointer",
				doReference(&ast.Ident{
					Name: name,
				}),
			),
		},
	})

	af.AddAssignment("cp_"+name, &ast.CallExpr{
		Fun: &ast.ParenExpr{
			X: doPointer(sliceType),
		},
		Args: []ast.Expr{
			doCall(
				"unsafe",
				"Pointer",
				accessMember("gos_"+name, "Data"),
			),
		},
	})
}

// AddGoToCSliceConversion adds Go to C slice conversion to af.
func (af *ASTFunc) AddGoToCSliceConversion(name string, typ Type) {
	// Declare the slice
	sliceType := sliceType(typ)

	af.AddAssignment(
		"ca_"+name,
		doCast(
			"make",
			&ast.ArrayType{
				Elt: sliceType,
			},
			doCast(
				"len",
				&ast.Ident{
					Name: name,
				},
			),
		),
	)

	af.AddStatement(doDeclare(
		"cp_"+name,
		doPointer(sliceType),
	))

	af.AddStatement(&ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X: doCast(
				"len",
				&ast.Ident{
					Name: name,
				},
			),
			Op: token.GTR,
			Y:  doZero(),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "cp_" + name,
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						doReference(&ast.IndexExpr{
							X: &ast.Ident{
								Name: "ca_" + name,
							},
							Index: doZero(),
						}),
					},
				},
			},
		},
	})

	// assign elements
	var loopStatements []ast.Stmt

	// handle our good old friend the const char * differently...
	if typ.CGoName == CSChar {
		loopStatements = append(loopStatements, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "ci_str",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				doCCast(
					"CString",
					&ast.IndexExpr{
						X: &ast.Ident{
							Name: name,
						},
						Index: &ast.Ident{
							Name: "i",
						},
					},
				),
			},
		})

		loopStatements = append(loopStatements, &ast.DeferStmt{
			Call: doCCast(
				"free",
				doCall(
					"unsafe",
					"Pointer",
					&ast.Ident{
						Name: "ci_str",
					},
				),
			),
		})

		loopStatements = append(loopStatements, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.IndexExpr{
					X: &ast.Ident{
						Name: "ca_" + name,
					},
					Index: &ast.Ident{
						Name: "i",
					},
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.Ident{
					Name: "ci_str",
				},
			},
		})

	} else {
		loopStatements = append(loopStatements, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.IndexExpr{
					X: &ast.Ident{
						Name: "ca_" + name,
					},
					Index: &ast.Ident{
						Name: "i",
					},
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.SelectorExpr{
					X: &ast.IndexExpr{
						X: &ast.Ident{
							Name: name,
						},
						Index: &ast.Ident{
							Name: "i",
						},
					},
					Sel: &ast.Ident{
						Name: "c",
					},
				},
			},
		})
	}

	af.AddStatement(&ast.RangeStmt{
		Key: &ast.Ident{
			Name: "i",
		},
		Tok: token.DEFINE,
		X: &ast.Ident{
			Name: name,
		},
		Body: &ast.BlockStmt{
			List: loopStatements,
		},
	})
}
