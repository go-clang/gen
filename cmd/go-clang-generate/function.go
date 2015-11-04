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
	fa := NewASTFunc(f)

	hasReturnArguments := false

	// TODO maybe name the return arguments ... because of clang_getDiagnosticOption -> the normal return can be always just "o"?

	// TODO reenable this, see the comment at the bottom of the generate function for details
	// Add function comment
	/*if f.Comment != "" {
		fa.Doc = &ast.CommentGroup{
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
			fa.Recv = &ast.FieldList{
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

		fa.Type.Params = &ast.FieldList{
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

				fa.addAssignment(
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
				fa.addStatement(doDeclare(
					"cp_"+p.Name,
					doPointer(sliceType),
				))
				fa.addStatement(&ast.IfStmt{
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
									doReference(&ast.IndexExpr{
										X: &ast.Ident{
											Name: "ca_" + p.Name,
										},
										Index: doZero(),
									}),
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

				fa.addStatement(&ast.RangeStmt{
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

					fa.addReturnType("", retType)
				}

				if p.Type.IsSlice && p.Type.IsReturnArgument {
					fa.addStatement(doDeclare(
						"cp_"+p.Name,
						getSliceType(p.Type),
					))
				}

				// Declare the return argument's variable
				var varType ast.Expr
				if p.Type.PointerLevel > 0 && p.Type.CGoName == CSChar {
					varType = doPointer(doCType("char"))
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

				fa.addStatement(doDeclare(
					p.Name,
					varType,
				))

				if p.Type.GoName == "cxstring" {
					fa.addDefer(doCall(p.Name, "Dispose"))
				} else if p.Type.PointerLevel > 0 && p.Type.CGoName == CSChar {
					fa.addDefer(doCCast(
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
						fa.addReturnItem(doCCast(
							"GoString",
							&ast.Ident{
								Name: p.Name,
							},
						))
					} else if p.Type.GoName == "cxstring" {
						fa.addReturnItem(doCall(p.Name, "String"))
					} else if p.Type.IsPrimitive {
						fa.addReturnItem(doCast(
							p.Type.GoName,
							&ast.Ident{
								Name: p.Name,
							},
						))
					} else {
						fa.addReturnItem(&ast.Ident{
							Name: p.Name,
						})
					}
				}

				continue
			} else if p.Type.PointerLevel > 0 && p.Type.IsPrimitive && p.Type.CGoName != CSChar {
				hasDeclaration = true

				var varType = doCType(p.Type.CGoName)

				fa.addStatement(doDeclare(
					"cp_"+p.Name,
					varType,
				))
				fa.addStatement(&ast.IfStmt{
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

			fa.Type.Params.List = append(fa.Type.Params.List, doField(p.Name, p.Type))
		}

		if hasReturnArguments || hasDeclaration {
			fa.addEmptyLine()
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

					fa.addAssignment(
						"c_"+p.Name,
						doCCast(
							"CString",
							&ast.Ident{
								Name: p.Name,
							},
						),
					)
					fa.addDefer(doCCast(
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
						pf = doReference(&ast.Ident{
							Name: "cp_" + p.Name,
						})
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
					pf = doReference(pf)
				}
			}

			if p.Type.IsReturnArgument {
				pf = doReference(pf)
			}

			call.Args = append(call.Args, pf)
		}

		if goToCTypeConversions {
			fa.addEmptyLine()
		}
	}

	// Check if we need to add a return
	if f.ReturnType.GoName != "void" || hasReturnArguments {
		if f.ReturnType.GoName == "cxstring" {
			// Do the C function call and save the result into the new variable "o" while transforming it into a cxstring
			fa.addAssignment("o", doCompose("cxstring", call))
			fa.addDefer(doCall("o", "Dispose"))
			fa.addEmptyLine()

			// Call the String method on the cxstring instance
			fa.addReturnItem(doCall("o", "String"))

			// Change the return type to "string"
			fa.addReturnType("", Type{
				GoName: "string",
			})
		} else {
			if f.ReturnType.GoName != "void" {
				// Add the function return type
				fa.addReturnType("", f.ReturnType)
			}

			// Do we need to convert the return of the C function into a boolean?
			if f.ReturnType.GoName == "bool" {
				// Do the C function call and save the result into the new variable "o"
				fa.addAssignment("o", call)
				fa.addEmptyLine()

				// Check if o is not equal to zero and return the result
				fa.addReturnItem(&ast.BinaryExpr{
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
				fa.addReturnItem(doCCast(
					"GoString",
					call,
				))
			} else if f.ReturnType.GoName == "time.Time" {
				fa.addReturnItem(doCall(
					"time",
					"Unix",
					doCast("int64", call),
					doZero(),
				))
			} else if f.ReturnType.GoName == "void" {
				// Handle the case where the C function has no return argument but parameters that are return arguments

				// Do the C function call
				fa.addStatement(&ast.ExprStmt{
					X: call,
				})
				fa.addEmptyLine()
			} else if f.ReturnType.PointerLevel > 0 {
				// Do the C function call and save the result into the new variable "o"
				fa.addAssignment("o", doUnreference(call))
				fa.addEmptyLine()

				fa.addReturnItem(doReference(doCompose(
					f.ReturnType.GoName,
					&ast.Ident{
						Name: "o",
					},
				)))
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
					fa.addAssignment("o", convCall)
					fa.addEmptyLine()

					// Add the C function call result to the return statement
					fa.addReturnItem(&ast.Ident{
						Name: "o",
					})
				} else {
					fa.addReturnItem(convCall)
				}
			}
		}

		fa.addCToGoConversions()

		// Add the return statement
		fa.addStatement(fa.Return)
	} else {
		fa.addCToGoConversions()

		// No return needed, just add the C function call
		fa.addStatement(&ast.ExprStmt{
			X: call,
		})
	}

	var b bytes.Buffer
	err := format.Node(&b, token.NewFileSet(), []ast.Decl{fa.FuncDecl})
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
