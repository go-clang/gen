package main

import (
	"go/ast"
	"go/token"
)

func accessMember(variable string, member string) *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X: &ast.Ident{
			Name: variable,
		},
		Sel: &ast.Ident{
			Name: member,
		},
	}
}

func doCast(typ string, args ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun: &ast.Ident{
			Name: typ,
		},
		Args: args,
	}
}

func doCompose(typ string, v ast.Expr) *ast.CompositeLit {
	return &ast.CompositeLit{
		Type: &ast.Ident{
			Name: typ,
		},
		Elts: []ast.Expr{
			v,
		},
	}
}

func doCType(c string) *ast.SelectorExpr {
	return accessMember("C", c)
}

func doCCast(typ string, args ...ast.Expr) *ast.CallExpr {
	return doCall("C", typ, args...)
}

func doDeclare(name string, typ ast.Expr) *ast.DeclStmt {
	return &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{
						&ast.Ident{
							Name: name,
						},
					},
					Type: typ,
				},
			},
		},
	}
}

func doField(name string, typ Type) *ast.Field {
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

func doCall(variable string, method string, args ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun:  accessMember(variable, method),
		Args: args,
	}
}

func doReference(x ast.Expr) *ast.UnaryExpr {
	return &ast.UnaryExpr{
		Op: token.AND,
		X:  x,
	}
}

func doPointer(x ast.Expr) *ast.StarExpr {
	return &ast.StarExpr{
		X: x,
	}
}

func doUnreference(x ast.Expr) *ast.StarExpr {
	return &ast.StarExpr{
		X: x,
	}
}

func doZero() *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: "0",
	}
}

func getSliceType(typ Type) ast.Expr {
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
