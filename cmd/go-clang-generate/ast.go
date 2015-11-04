package main

import (
	"go/ast"
	"go/token"
)

type ASTFunc struct {
	*ast.FuncDecl

	Function *Function
	Return   *ast.ReturnStmt
}

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

		Function: f,
		Return: &ast.ReturnStmt{
			Results: []ast.Expr{},
		},
	}
}

func (fa *ASTFunc) addStatement(stmt ast.Stmt) {
	fa.Body.List = append(fa.Body.List, stmt)
}

func (fa *ASTFunc) addAssignment(variable string, e ast.Expr) {
	fa.addStatement(&ast.AssignStmt{
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

func (fa *ASTFunc) addDefer(call *ast.CallExpr) {
	fa.addStatement(&ast.DeferStmt{
		Call: call,
	})
}

func (fa *ASTFunc) addEmptyLine() {
	// TODO this should be done using something else than a fake statement.
	fa.addStatement(&ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.Ident{
				Name: "REMOVE",
			},
		},
	})
}

func (fa *ASTFunc) addReturnItem(item ast.Expr) {
	fa.Return.Results = append(fa.Return.Results, item)
}

func (fa *ASTFunc) addReturnType(name string, typ Type) {
	fa.Type.Results.List = append(fa.Type.Results.List, doField(name, typ))
}

func (fa *ASTFunc) addCToGoConversions() {
	cToGoTypeConversions := false

	for _, p := range fa.Function.Parameters {
		if p.Type.IsSlice && p.Type.IsReturnArgument {
			cToGoTypeConversions = true

			var lengthOfSlice string
			for _, pl := range fa.Function.Parameters {
				if pl.Type.LengthOfSlice == p.Name {
					lengthOfSlice = pl.Name

					break
				}
			}

			fa.addAssignment("gos_"+p.Name, &ast.CallExpr{
				Fun: &ast.ParenExpr{
					X: doPointer(accessMember("reflect", "SliceHeader")),
				},
				Args: []ast.Expr{
					doCall(
						"unsafe",
						"Pointer",
						doReference(&ast.Ident{
							Name: p.Name,
						}),
					),
				},
			})
			fa.addStatement(&ast.AssignStmt{
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
			fa.addStatement(&ast.AssignStmt{
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
			fa.addStatement(&ast.AssignStmt{
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
		fa.addEmptyLine()
	}
}
