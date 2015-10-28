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

type Function struct {
	Name    string
	CName   string
	Comment string

	Parameters []FunctionParameter
	ReturnType Type

	Receiver Receiver

	Member string
}

type FunctionParameter struct {
	Name  string
	CName string
	Type  Type
}

func handleFunctionCursor(cursor clang.Cursor) *Function {
	f := Function{
		CName:   cursor.Spelling(),
		Comment: cleanDoxygenComment(cursor.RawCommentText()),

		Parameters: []FunctionParameter{},
		ReturnType: Type{
			Name: trimClangPrefix(cursor.ResultType().TypeSpelling()),
		},
	}

	f.Name = strings.TrimPrefix(f.CName, "clang_")

	numParam := uint(cursor.NumArguments())
	for i := uint(0); i < numParam; i++ {
		param := cursor.Argument(i)

		p := FunctionParameter{
			CName: param.DisplayName(),
			Type: Type{
				CName: param.Type().TypeSpelling(),
			},
		}

		p.Name = p.CName
		p.Type.Name = trimClangPrefix(p.Type.CName)

		if p.Name == "" {
			p.Name = receiverName(p.Type.Name)
		}

		f.Parameters = append(f.Parameters, p)
	}

	return &f
}

func generateASTFunction(f *Function) string {
	astFunc := ast.FuncDecl{
		Name: &ast.Ident{
			Name: f.Name,
		},
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{},
	}

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
					&ast.Field{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: f.Receiver.Name,
							},
						},
						Type: &ast.Ident{
							Name: f.Receiver.Type.Name,
						},
					},
				},
			}
		}
	}

	// Basic call to the C function
	call := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "C",
			},
			Sel: &ast.Ident{
				Name: f.CName,
			},
		},
		Args: []ast.Expr{},
	}

	if len(f.Parameters) != 0 {
		if f.Receiver.Name != "" {
			f.Parameters[0].Name = f.Receiver.Name
		}

		astFunc.Type.Params = &ast.FieldList{
			List: []*ast.Field{},
		}

		// Add parameters to the function
		for i, p := range f.Parameters {
			if i == 0 && f.Receiver.Name != "" {
				continue
			}

			astFunc.Type.Params.List = append(astFunc.Type.Params.List, &ast.Field{
				Names: []*ast.Ident{
					&ast.Ident{
						Name: p.Name,
					},
				},
				Type: &ast.Ident{
					Name: p.Type.Name,
				},
			})
		}

		goToCTypeConversions := false

		// Add arguments to the C function call
		for _, p := range f.Parameters {
			if p.Type.Primitive != "" {
				// Handle Go type to C type conversions
				if p.Type.Primitive == "const char *" {
					goToCTypeConversions = true

					astFunc.Body.List = append(astFunc.Body.List, &ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: "c_" + p.Name,
							},
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "C",
									},
									Sel: &ast.Ident{
										Name: "CString",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: p.Name,
									},
								},
							},
						},
					})
					astFunc.Body.List = append(astFunc.Body.List, &ast.DeferStmt{
						Call: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "C",
								},
								Sel: &ast.Ident{
									Name: "free",
								},
							},
							Args: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "unsafe",
										},
										Sel: &ast.Ident{
											Name: "Pointer",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "c_" + p.Name,
										},
									},
								},
							},
						},
					})

					call.Args = append(call.Args, &ast.Ident{
						Name: "c_" + p.Name,
					})
				} else if p.Type.Primitive == "cxstring" { // TODO try to get cxstring and "String" completely out of this function since it is just a struct which can be handled by the struct code
					call.Args = append(call.Args, &ast.SelectorExpr{
						X: &ast.Ident{
							Name: p.Name,
						},
						Sel: &ast.Ident{
							Name: "c",
						},
					})
				} else {
					call.Args = append(call.Args, &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "C",
							},
							Sel: &ast.Ident{
								Name: p.Type.Primitive,
							},
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: p.Name,
							},
						},
					})
				}
			} else {
				call.Args = append(call.Args, &ast.SelectorExpr{
					X: &ast.Ident{
						Name: p.Name,
					},
					Sel: &ast.Ident{
						Name: "c",
					},
				})
			}
		}

		if goToCTypeConversions {
			// TODO maybe somehow remove this?! We add an empty line here
			astFunc.Body.List = append(astFunc.Body.List, &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.Ident{
						Name: "REMOVE",
					},
				},
			})
		}
	}

	// Check if we need to add a return
	if f.ReturnType.Name != "void" {
		// Add the function return type
		astFunc.Type.Results = &ast.FieldList{
			List: []*ast.Field{
				&ast.Field{
					Type: &ast.Ident{
						Name: f.ReturnType.Name,
					},
				},
			},
		}

		// Convert the return value of the C function
		var convCall ast.Expr

		// Structs are literals, everything else is a cast
		if f.ReturnType.Primitive == "" {
			convCall = &ast.CompositeLit{
				Type: &ast.Ident{
					Name: f.ReturnType.Name,
				},
				Elts: []ast.Expr{
					call,
				},
			}
		} else {
			convCall = &ast.CallExpr{
				Fun: &ast.Ident{
					Name: f.ReturnType.Name,
				},
				Args: []ast.Expr{
					call,
				},
			}
		}

		result := convCall

		// Do we need to convert the return of the C function into a boolean?
		if f.ReturnType.Name == "bool" && f.ReturnType.Primitive != "" {
			// Do the C function call and save the result into the new variable "o"
			astFunc.Body.List = append(astFunc.Body.List, &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: "o",
					},
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					call, // No cast needed
				},
			})

			// TODO maybe somehow remove this?! We add an empty line here
			astFunc.Body.List = append(astFunc.Body.List, &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.Ident{
						Name: "REMOVE",
					},
				},
			})

			// Check if o is not equal to zero and return the result
			result = &ast.BinaryExpr{
				X: &ast.Ident{
					Name: "o",
				},
				Op: token.NEQ,
				Y: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "C",
						},
						Sel: &ast.Ident{
							Name: f.ReturnType.Primitive,
						},
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.INT,
							Value: "0",
						},
					},
				},
			}
		} else if f.ReturnType.Name == "string" {
			// If this is a normal const char * C type there is not so much to do
			if f.ReturnType.Primitive == "const char *" {
				result = &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "C",
						},
						Sel: &ast.Ident{
							Name: "GoString",
						},
					},
					Args: []ast.Expr{
						call,
					},
				}
			} else {
				// This should be a CXString so handle it accordingly

				// Do the C function call and save the result into the new variable "o" while transforming it into a cxstring
				astFunc.Body.List = append(astFunc.Body.List, &ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "o",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.Ident{
								Name: "cxstring",
							},
							Elts: []ast.Expr{
								call,
							},
						},
					},
				})
				astFunc.Body.List = append(astFunc.Body.List, &ast.DeferStmt{
					Call: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "o",
							},
							Sel: &ast.Ident{
								Name: "Dispose",
							},
						},
					},
				})

				// TODO maybe somehow remove this?! We add an empty line here
				astFunc.Body.List = append(astFunc.Body.List, &ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "REMOVE",
						},
					},
				})

				// Call the String method on the cxstring instance
				result = &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "o",
						},
						Sel: &ast.Ident{
							Name: "String",
						},
					},
				}
			}
		} else if f.ReturnType.Name == "time.Time" {
			result = &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "time",
					},
					Sel: &ast.Ident{
						Name: "Unix",
					},
				},
				Args: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.Ident{
							Name: "int64",
						},
						Args: []ast.Expr{
							call,
						},
					},
					&ast.BasicLit{
						Kind:  token.INT,
						Value: "0",
					},
				},
			}
		}

		// Add the return statement
		astFunc.Body.List = append(astFunc.Body.List, &ast.ReturnStmt{
			Results: []ast.Expr{
				result,
			},
		})
	} else {
		// No return needed, just add the C function call
		astFunc.Body.List = append(astFunc.Body.List, &ast.ExprStmt{
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
func ({{$.Receiver}} {{$.ReceiverType}}) {{$.Name}}() {{if $.Type.PointerLevel}}*{{end}}{{if $.ReturnType.Primitive}}{{$.ReturnType.Primitive}}{{else}}{{$.ReturnType}}{{end}} {
	return {{if $.Type.PointerLevel}}&{{end}}{{if $.ReturnType.Primitive}}{{$.ReturnType.Primitive}}{{else}}{{$.ReturnType}}{{end}}{{if $.ReturnType.Primitive}}({{if $.Type.PointerLevel}}*{{end}}{{$.Receiver}}.c.{{$.Member}}){{else}}{{"{"}}{{if $.Type.PointerLevel}}*{{end}}{{$.Receiver}}.c.{{$.Member}}{{"}"}}{{end}}
}
`))

func generateFunctionStructMemberGetter(f *Function) string {
	var b bytes.Buffer
	if err := templateGenerateStructMemberGetter.Execute(&b, f); err != nil {
		panic(err)
	}

	return b.String()
}

type FunctionSliceReturn struct {
	Function

	ElementType     string
	IsPrimitive     bool
	ArrayDimensions int
}

var templateGenerateReturnSlice = template.Must(template.New("go-clang-generate-slice").Parse(`{{$.Comment}}
func ({{$.Receiver}} {{$.ReceiverType}}) {{$.Name}}() []{{if eq $.ArrayDimensions 2 }}*{{end}}{{$.ElementType}} {
	sc := []{{$.ElementType}}{}
	{{if eq $.ArrayDimensions 2 }}
	length := int(C.sizeof({{$.Receiver}}.c.{{$.Member}}[0])) / int(sizeof({{$.Receiver}}.c.{{$.Member}}[0][0]))
	{{else}}
	length := int(sizeof({{$.Receiver}}.c.{{$.Member}}))
	{{end}}
	for is := 0; is < length; is++ {
		sc = append(sc, {{if eq $.ArrayDimensions 2}}&{{$.ElementType}}{{else}}{{$.ElementType}}{{end}}{{if $.IsPrimitive}}({{$.Receiver}}.c.{{$.Member}}[is])){{else}}{"{"}{{$.Receiver}}.c.{{$.Member}}[is]){{"}"}}{{end}}
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

func generateFunction(name, cname, comment, member string, typ Type) *Function {
	receiverType := trimClangPrefix(cname)
	receiverName := receiverName(receiverType)
	functionName := upperFirstCharacter(name)

	rType := ""
	rTypePrimitive := ""

	if typ.IsPrimitive {
		rTypePrimitive = typ.Name
	} else {
		rType = typ.Name
	}

	f := &Function{
		Name:    functionName,
		CName:   cname,
		Comment: comment,

		Parameters: []FunctionParameter{},

		ReturnType: Type{
			Name:      rType,
			Primitive: rTypePrimitive,

			PointerLevel:  typ.PointerLevel,
			IsEnumLiteral: typ.IsEnumLiteral,
		},

		Receiver: Receiver{
			Name: receiverName,
			Type: Type{
				Name: receiverType,
			},
		},

		Member: member,
	}

	return f
}
