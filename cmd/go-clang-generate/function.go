package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
	"text/template"

	"github.com/sbinet/go-clang"
)

func trimCommonFunctionName(name string, typ Type) string {
	name = trimCommonFunctionNamePrefix(name)

	if fn := strings.TrimPrefix(name, typ.GoName+"_"); len(fn) != len(name) {
		name = fn
	} else if fn := strings.TrimPrefix(name, typ.GoName); len(fn) != len(name) {
		name = fn
	}

	name = trimCommonFunctionNamePrefix(name)

	// If the function name is empty at this point, it is a constructor
	if name == "" {
		name = typ.GoName
	}

	return name
}

func trimCommonFunctionNamePrefix(name string) string {
	name = strings.TrimPrefix(name, "create")
	name = strings.TrimPrefix(name, "get")

	name = trimLanguagePrefix(name)

	return name
}

func trimLanguagePrefix(name string) string {
	name = strings.TrimPrefix(name, "CX_CXX")
	name = strings.TrimPrefix(name, "CXX")
	name = strings.TrimPrefix(name, "CX")
	name = strings.TrimPrefix(name, "ObjC")

	return name
}

// Function represents a generation function
type Function struct {
	Name    string
	CName   string
	Comment string

	Parameters []FunctionParameter
	ReturnType Type

	Receiver Receiver

	Member string
}

// FunctionParameter represents a generation function parameter
type FunctionParameter struct {
	Name  string
	CName string
	Type  Type
}

func NewFunction(name, cname, comment, member string, typ Type) *Function {
	receiverType := trimLanguagePrefix(cname)
	receiverName := receiverName(receiverType)
	functionName := upperFirstCharacter(name)

	if typ.IsPrimitive {
		typ.CGoName = typ.GoName
	}
	if (strings.HasPrefix(name, "has") || strings.HasPrefix(name, "is")) && typ.GoName == GoInt16 {
		typ.GoName = GoBool
	}

	f := &Function{
		Name:    functionName,
		CName:   cname,
		Comment: comment,

		Parameters: []FunctionParameter{},

		ReturnType: typ,
		Receiver: Receiver{
			Name: receiverName,
			Type: Type{
				GoName: receiverType,
			},
		},

		Member: member,
	}

	return f
}

func HandleFunctionCursor(cursor clang.Cursor) *Function {
	f := Function{
		CName:   cursor.Spelling(),
		Comment: cleanDoxygenComment(cursor.RawCommentText()),

		Parameters: []FunctionParameter{},
	}

	f.Name = strings.TrimPrefix(f.CName, "clang_")

	typ, err := TypeFromClangType(cursor.ResultType())
	if err != nil {
		panic(err)
	}
	f.ReturnType = typ

	numParam := uint(cursor.NumArguments())
	for i := uint(0); i < numParam; i++ {
		param := cursor.Argument(i)

		p := FunctionParameter{
			CName: param.DisplayName(),
		}

		typ, err := TypeFromClangType(param.Type())
		if err != nil {
			panic(err)
		}
		p.Type = typ

		p.Name = p.CName
		if p.Name == "" {
			p.Name = receiverName(p.Type.GoName)
		} else {
			pns := strings.Split(p.Name, "_")
			for i := range pns {
				pns[i] = upperFirstCharacter(pns[i])
			}
			p.Name = lowerFirstCharacter(strings.Join(pns, ""))
		}
		if r := ReplaceGoKeywords(p.Name); r != "" {
			p.Name = r
		}

		f.Parameters = append(f.Parameters, p)
	}

	return &f
}

func (f *Function) Generate() string {
	astFunc := ast.FuncDecl{
		Name: &ast.Ident{
			Name: f.Name,
		},
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{},
			},
		},
		Body: &ast.BlockStmt{},
	}

	retur := &ast.ReturnStmt{
		Results: []ast.Expr{},
	}
	hasReturnArguments := false

	accessMember := func(variable string, member string) *ast.SelectorExpr {
		return &ast.SelectorExpr{
			X: &ast.Ident{
				Name: variable,
			},
			Sel: &ast.Ident{
				Name: member,
			},
		}
	}
	addStatement := func(stmt ast.Stmt) {
		astFunc.Body.List = append(astFunc.Body.List, stmt)
	}
	addAssignment := func(variable string, e ast.Expr) {
		addStatement(&ast.AssignStmt{
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
	addAssignmentToO := func(e ast.Expr) {
		addAssignment("o", e)
	}
	addDefer := func(call *ast.CallExpr) {
		addStatement(&ast.DeferStmt{
			Call: call,
		})
	}
	addEmptyLine := func() {
		// TODO this should be done using something else than a fake statement.
		addStatement(&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.Ident{
					Name: "REMOVE",
				},
			},
		})
	}
	addReturnItem := func(item ast.Expr) {
		retur.Results = append(retur.Results, item)
	}
	doCall := func(variable string, method string, args ...ast.Expr) *ast.CallExpr {
		return &ast.CallExpr{
			Fun:  accessMember(variable, method),
			Args: args,
		}
	}
	doCast := func(typ string, args ...ast.Expr) *ast.CallExpr {
		return &ast.CallExpr{
			Fun: &ast.Ident{
				Name: typ,
			},
			Args: args,
		}
	}
	doCompose := func(typ string, v ast.Expr) *ast.CompositeLit {
		return &ast.CompositeLit{
			Type: &ast.Ident{
				Name: typ,
			},
			Elts: []ast.Expr{
				v,
			},
		}
	}
	doCType := func(c string) *ast.SelectorExpr {
		return accessMember("C", c)
	}
	doCCast := func(typ string, args ...ast.Expr) *ast.CallExpr {
		return doCall("C", typ, args...)
	}
	doField := func(name string, typ Type) *ast.Field {
		f := &ast.Field{}

		if name != "" {
			f.Names = []*ast.Ident{
				&ast.Ident{
					Name: name,
				},
			}
		}
		if typ.GoName != "" {
			if typ.PointerLevel > 0 && typ.CGoName == CSChar {
				f.Type = &ast.Ident{
					Name: "string",
				}
			} else {
				f.Type = &ast.Ident{
					Name: typ.GoName,
				}
			}

			if typ.IsSlice {
				f.Type = &ast.ArrayType{
					Elt: f.Type,
				}
			} else if typ.PointerLevel > 0 && typ.CGoName != CSChar && !typ.IsReturnArgument {
				for i := 0; i < typ.PointerLevel; i++ {
					f.Type = &ast.StarExpr{
						X: f.Type,
					}
				}
			}
		}

		return f
	}
	addReturnType := func(name string, typ Type) {
		astFunc.Type.Results.List = append(astFunc.Type.Results.List, doField(name, typ))
	}
	doZero := func() *ast.BasicLit {
		return &ast.BasicLit{
			Kind:  token.INT,
			Value: "0",
		}
	}
	getSliceType := func(typ Type) ast.Expr {
		var sliceType ast.Expr

		if typ.PointerLevel > 0 && typ.CGoName == CSChar {
			sliceType = doCType("char")
		} else {
			sliceType = doCType(typ.CGoName)
		}

		for i := 1; i < typ.PointerLevel; i++ {
			sliceType = &ast.StarExpr{
				X: sliceType,
			}
		}

		return sliceType
	}

	addCToGoConversions := func() {
		cToGoTypeConversions := false

		for _, p := range f.Parameters {
			if p.Type.IsSlice && p.Type.IsReturnArgument {
				cToGoTypeConversions = true

				var lengthOfSlice string
				for _, pl := range f.Parameters {
					if pl.Type.LengthOfSlice == p.Name {
						lengthOfSlice = pl.Name

						break
					}
				}

				addAssignment("gos_"+p.Name, &ast.CallExpr{
					Fun: &ast.ParenExpr{
						X: &ast.StarExpr{
							X: accessMember("reflect", "SliceHeader"),
						},
					},
					Args: []ast.Expr{
						doCall(
							"unsafe",
							"Pointer",
							&ast.UnaryExpr{
								Op: token.AND,
								X: &ast.Ident{
									Name: p.Name,
								},
							},
						),
					},
				})
				addStatement(&ast.AssignStmt{
					Lhs: []ast.Expr{
						accessMember("gos_"+p.Name, "Cap"),
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
				addStatement(&ast.AssignStmt{
					Lhs: []ast.Expr{
						accessMember("gos_"+p.Name, "Len"),
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
				addStatement(&ast.AssignStmt{
					Lhs: []ast.Expr{
						accessMember("gos_"+p.Name, "Data"),
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						doCast(
							"uintptr",
							doCall(
								"unsafe",
								"Pointer",
								&ast.Ident{
									Name: "cp_" + p.Name,
								},
							),
						),
					},
				})
			}
		}

		if cToGoTypeConversions {
			addEmptyLine()
		}
	}

	// TODO maybe name the return arguments ... because of clang_getDiagnosticOption -> the normal return can be always just "o"?

	// TODO reenable this, see the comment at the bottom of the generate function for details
	// Add function comment
	/*if f.Comment != "" {
		astFunc.Doc = &ast.CommentGroup{
			List: []*ast.Comment{
				&ast.Comment{
					Text: f.Comment,
				},
			},
		}
	}*/

	// Add receiver to make function a method
	if f.Receiver.Name != "" {
		if len(f.Parameters) > 0 { // TODO maybe to not set the receiver at all? -> do this outside of the generate function?
			astFunc.Recv = &ast.FieldList{
				List: []*ast.Field{
					doField(f.Receiver.Name, f.Receiver.Type),
				},
			}
		}
	}

	// Basic call to the C function
	call := doCCast(f.CName)

	if len(f.Parameters) != 0 {
		if f.Receiver.Name != "" {
			f.Parameters[0].Name = f.Receiver.Name
		}

		astFunc.Type.Params = &ast.FieldList{
			List: []*ast.Field{},
		}

		hasDeclaration := false

		// Add parameters to the function
		for i, p := range f.Parameters {
			if i == 0 && f.Receiver.Name != "" {
				continue
			}

			// Ingore length parameters since they will be filled by the slice itself
			if p.Type.LengthOfSlice != "" && !p.Type.IsReturnArgument {
				continue
			}

			if p.Type.IsSlice && !p.Type.IsReturnArgument {
				hasDeclaration = true

				// Declare the slice
				sliceType := getSliceType(p.Type)

				addAssignment(
					"ca_"+p.Name,
					doCast(
						"make",
						&ast.ArrayType{
							Elt: sliceType,
						},
						doCast(
							"len",
							&ast.Ident{
								Name: p.Name,
							},
						),
					),
				)
				addStatement(&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "cp_" + p.Name,
									},
								},
								Type: &ast.StarExpr{
									X: sliceType,
								},
							},
						},
					},
				})
				addStatement(&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: doCast(
							"len",
							&ast.Ident{
								Name: p.Name,
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
										Name: "cp_" + p.Name,
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.UnaryExpr{
										Op: token.AND,
										X: &ast.IndexExpr{
											X: &ast.Ident{
												Name: "ca_" + p.Name,
											},
											Index: doZero(),
										},
									},
								},
							},
						},
					},
				})

				// Assign elements
				var loopStatements []ast.Stmt

				// Handle our good old friend the const char * differently...
				if p.Type.CGoName == CSChar {
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
										Name: p.Name,
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
									Name: "ca_" + p.Name,
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
									Name: "ca_" + p.Name,
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
										Name: p.Name,
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

				addStatement(&ast.RangeStmt{
					Key: &ast.Ident{
						Name: "i",
					},
					Tok: token.DEFINE,
					X: &ast.Ident{
						Name: p.Name,
					},
					Body: &ast.BlockStmt{
						List: loopStatements,
					},
				})
			} else if p.Type.IsReturnArgument {
				hasReturnArguments = true

				if p.Type.LengthOfSlice == "" {
					// Add the return type to the function return arguments
					retType := p.Type
					if p.Type.GoName == "cxstring" {
						retType.GoName = "string"
					}

					addReturnType("", retType)
				}

				if p.Type.IsSlice && p.Type.IsReturnArgument {
					addStatement(&ast.DeclStmt{
						Decl: &ast.GenDecl{
							Tok: token.VAR,
							Specs: []ast.Spec{
								&ast.ValueSpec{
									Names: []*ast.Ident{
										&ast.Ident{
											Name: "cp_" + p.Name,
										},
									},
									Type: getSliceType(p.Type),
								},
							},
						},
					})
				}

				// Declare the return argument's variable
				var varType ast.Expr
				if p.Type.PointerLevel > 0 && p.Type.CGoName == CSChar {
					varType = &ast.StarExpr{
						X: doCType("char"),
					}
				} else if p.Type.IsPrimitive {
					varType = doCType(p.Type.CGoName)
				} else {
					varType = &ast.Ident{
						Name: p.Type.GoName,
					}
				}
				if p.Type.IsSlice {
					varType = &ast.ArrayType{
						Elt: varType,
					}
				}

				addStatement(&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: p.Name,
									},
								},
								Type: varType,
							},
						},
					},
				})
				if p.Type.GoName == "cxstring" {
					addDefer(doCall(p.Name, "Dispose"))
				} else if p.Type.PointerLevel > 0 && p.Type.CGoName == CSChar {
					addDefer(doCCast(
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
					// Add the return argument to the return statement
					if p.Type.PointerLevel > 0 && p.Type.CGoName == CSChar {
						addReturnItem(doCCast(
							"GoString",
							&ast.Ident{
								Name: p.Name,
							},
						))
					} else if p.Type.GoName == "cxstring" {
						addReturnItem(doCall(p.Name, "String"))
					} else if p.Type.IsPrimitive {
						addReturnItem(doCast(
							p.Type.GoName,
							&ast.Ident{
								Name: p.Name,
							},
						))
					} else {
						addReturnItem(&ast.Ident{
							Name: p.Name,
						})
					}
				}

				continue
			} else if p.Type.PointerLevel > 0 && p.Type.IsPrimitive && p.Type.CGoName != CSChar {
				hasDeclaration = true

				var varType = doCType(p.Type.CGoName)

				addStatement(&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "cp_" + p.Name,
									},
								},
								Type: varType,
							},
						},
					},
				})
				addStatement(&ast.IfStmt{
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
											&ast.StarExpr{
												X: &ast.Ident{
													Name: p.Name,
												},
											},
										},
									},
								},
							},
						},
					},
				})
			}

			astFunc.Type.Params.List = append(astFunc.Type.Params.List, doField(p.Name, p.Type))
		}

		if hasReturnArguments || hasDeclaration {
			addEmptyLine()
		}

		goToCTypeConversions := false

		// Add arguments to the C function call
		for _, p := range f.Parameters {
			var pf ast.Expr

			if p.Type.IsSlice {
				pf = &ast.Ident{
					Name: "cp_" + p.Name,
				}
			} else if p.Type.IsPrimitive {
				// Handle Go type to C type conversions
				if p.Type.PointerLevel == 1 && p.Type.CGoName == CSChar {
					goToCTypeConversions = true

					addAssignment(
						"c_"+p.Name,
						doCCast(
							"CString",
							&ast.Ident{
								Name: p.Name,
							},
						),
					)
					addDefer(doCCast(
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
				} else if p.Type.CGoName == "cxstring" { // TODO try to get cxstring and "String" completely out of this function since it is just a struct which can be handled by the struct code
					pf = accessMember(p.Name, "c")
				} else {
					if p.Type.IsReturnArgument {
						// Return arguments already have a cast
						pf = &ast.Ident{
							Name: p.Name,
						}
					} else if p.Type.LengthOfSlice != "" {
						pf = doCCast(
							p.Type.CGoName,
							doCast(
								"len",
								&ast.Ident{
									Name: p.Type.LengthOfSlice,
								},
							),
						)
					} else if p.Type.PointerLevel > 0 {
						pf = &ast.UnaryExpr{
							Op: token.AND,
							X: &ast.Ident{
								Name: "cp_" + p.Name,
							},
						}
					} else {
						pf = doCCast(
							p.Type.CGoName,
							&ast.Ident{
								Name: p.Name,
							},
						)
					}
				}
			} else {
				pf = accessMember(p.Name, "c")

				if p.Type.PointerLevel > 0 && !p.Type.IsReturnArgument {
					pf = &ast.UnaryExpr{
						Op: token.AND,
						X:  pf,
					}
				}
			}

			if p.Type.IsReturnArgument {
				pf = &ast.UnaryExpr{
					Op: token.AND,
					X:  pf,
				}
			}

			call.Args = append(call.Args, pf)
		}

		if goToCTypeConversions {
			addEmptyLine()
		}
	}

	// Check if we need to add a return
	if f.ReturnType.GoName != "void" || hasReturnArguments {
		if f.ReturnType.GoName == "cxstring" {
			// Do the C function call and save the result into the new variable "o" while transforming it into a cxstring
			addAssignmentToO(doCompose("cxstring", call))
			addDefer(doCall("o", "Dispose"))
			addEmptyLine()

			// Call the String method on the cxstring instance
			addReturnItem(doCall("o", "String"))

			// Change the return type to "string"
			addReturnType("", Type{
				GoName: "string",
			})
		} else {
			if f.ReturnType.GoName != "void" {
				// Add the function return type
				addReturnType("", f.ReturnType)
			}

			// Do we need to convert the return of the C function into a boolean?
			if f.ReturnType.GoName == "bool" {
				// Do the C function call and save the result into the new variable "o"
				addAssignmentToO(call)
				addEmptyLine()

				// Check if o is not equal to zero and return the result
				addReturnItem(&ast.BinaryExpr{
					X: &ast.Ident{
						Name: "o",
					},
					Op: token.NEQ,
					Y: doCCast(
						f.ReturnType.CGoName,
						doZero(),
					),
				})
			} else if f.ReturnType.CGoName == CSChar && f.ReturnType.PointerLevel == 1 {
				// If this is a normal const char * C type there is not so much to do
				addReturnItem(doCCast(
					"GoString",
					call,
				))
			} else if f.ReturnType.GoName == "time.Time" {
				addReturnItem(doCall(
					"time",
					"Unix",
					doCast("int64", call),
					doZero(),
				))
			} else if f.ReturnType.GoName == "void" {
				// Handle the case where the C function has no return argument but parameters that are return arguments

				// Do the C function call
				addStatement(&ast.ExprStmt{
					X: call,
				})
				addEmptyLine()
			} else if f.ReturnType.PointerLevel > 0 {
				// Do the C function call and save the result into the new variable "o"
				addAssignmentToO(&ast.StarExpr{
					X: call,
				})
				addEmptyLine()

				addReturnItem(&ast.UnaryExpr{
					Op: token.AND,
					X: doCompose(
						f.ReturnType.GoName,
						&ast.Ident{
							Name: "o",
						},
					),
				})
			} else {
				var convCall ast.Expr

				// Structs are literals, everything else is a cast
				if !f.ReturnType.IsPrimitive {
					convCall = doCompose(f.ReturnType.GoName, call)
				} else {
					convCall = doCast(f.ReturnType.GoName, call)
				}

				if hasReturnArguments {
					// Do the C function call and save the result into the new variable "o"
					addAssignmentToO(convCall)
					addEmptyLine()

					// Add the C function call result to the return statement
					addReturnItem(&ast.Ident{
						Name: "o",
					})
				} else {
					addReturnItem(convCall)
				}
			}
		}

		addCToGoConversions()

		// Add the return statement
		addStatement(retur)
	} else {
		addCToGoConversions()

		// No return needed, just add the C function call
		addStatement(&ast.ExprStmt{
			X: call,
		})
	}

	var b bytes.Buffer
	err := format.Node(&b, token.NewFileSet(), []ast.Decl{&astFunc})
	if err != nil {
		panic(err)
	}

	sss := b.String()

	// TODO hack to make new lines...
	sss = strings.Replace(sss, "REMOVE()", "", -1)

	// TODO find out how to position the comment correctly and do this using the AST
	if f.Comment != "" {
		sss = f.Comment + "\n" + sss
	}

	return sss
}

var templateGenerateStructMemberGetter = template.Must(template.New("go-clang-generate-function-getter").Parse(`{{$.Comment}}
func ({{$.Receiver.Name}} {{$.Receiver.Type.GoName}}) {{$.Name}}() {{if ge $.ReturnType.PointerLevel 1}}*{{end}}{{$.ReturnType.GoName}} {
	value := {{if eq $.ReturnType.GoName "bool"}}{{$.Receiver.Name}}.c.{{$.Member}}{{else}}{{$.ReturnType.GoName}}{{if $.ReturnType.IsPrimitive}}({{if ge $.ReturnType.PointerLevel 1}}*{{end}}{{$.Receiver.Name}}.c.{{$.Member}}){{else}}{{"{"}}{{if ge $.ReturnType.PointerLevel 1}}*{{end}}{{$.Receiver.Name}}.c.{{$.Member}}{{"}"}}{{end}}{{end}}
	return {{if eq $.ReturnType.GoName "bool"}}value != C.int(0){{else}}{{if ge $.ReturnType.PointerLevel 1}}&{{end}}value{{end}}
}
`))

func generateFunctionStructMemberGetter(f *Function) string {
	var b bytes.Buffer
	if err := templateGenerateStructMemberGetter.Execute(&b, f); err != nil {
		panic(err)
	}

	return b.String()
}

// FunctionSliceReturn TODO refactor
type FunctionSliceReturn struct {
	*Function

	SizeMember string

	CElementType    string
	ElementType     string
	IsPrimitive     bool
	ArrayDimensions int
	ArraySize       int64
}

var templateGenerateReturnSlice = template.Must(template.New("go-clang-generate-slice").Parse(`{{$.Comment}}
func ({{$.Receiver.Name}} {{$.Receiver.Type.GoName}}) {{$.Name}}() []{{if eq $.ArrayDimensions 2 }}*{{end}}{{$.ElementType}} {
	sc := []{{if eq $.ArrayDimensions 2 }}*{{end}}{{$.ElementType}}{}

	length := {{if ne $.ArraySize -1}}{{$.ArraySize}}{{else}}int({{$.Receiver.Name}}.c.{{$.SizeMember}}){{end}}
	goslice := (*[1 << 30]{{if or (eq $.ArrayDimensions 2) (eq $.ElementType "unsafe.Pointer")}}*{{end}}C.{{$.CElementType}})(unsafe.Pointer(&{{$.Receiver.Name}}.c.{{$.Member}}))[:length:length]

	for is := 0; is < length; is++ {
		sc = append(sc, {{if eq $.ArrayDimensions 2}}&{{end}}{{$.ElementType}}{{if $.IsPrimitive}}({{if eq $.ArrayDimensions 2}}*{{end}}goslice[is]){{else}}{{"{"}}{{if eq $.ArrayDimensions 2}}*{{end}}goslice[is]{{"}"}}{{end}})
	}

	return sc
}
`))

func generateFunctionSliceReturn(f *FunctionSliceReturn) string {
	var b bytes.Buffer
	if err := templateGenerateReturnSlice.Execute(&b, f); err != nil {
		panic(err)
	}

	return b.String()

}
