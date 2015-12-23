package expectations

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

func BuildUpExpectations(folder string) map[string][]string {
	expectations := make(map[string][]string)

	files, _ := ioutil.ReadDir(folder)
	for _, f := range files {
		structName := strings.Replace(f.Name(), ".go", "", 1)
		expectations[structName] = parseExpectations(fmt.Sprintf("%s%s", folder, f.Name()))
	}

	return expectations
}

func parseExpectations(file string) []string {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		panic(err)
	}

	expectations := []string{}

	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ValueSpec:
			for idx, ident := range x.Names {
				if strings.HasPrefix(ident.Name, "Test") {
					value := x.Values[idx]
					if _, ok := value.(*ast.BasicLit); ok && value.(*ast.BasicLit).Kind == token.STRING {
						exp := strings.Replace(strings.Replace(value.(*ast.BasicLit).Value, "`\n", "", -1), "\n`", "", -1)
						expectations = append(expectations, exp)
					}
				}
			}
		}
		return true
	})

	return expectations
}

func VerifyExpectations(fileContent string, expectations []string) bool {
	met := false
	for _, expct := range expectations {
		if !strings.Contains(fileContent, expct) {
			fmt.Printf("Expectation not met:\n %s \n", expct)
			met = met || true
		}
	}

	return !met
}
