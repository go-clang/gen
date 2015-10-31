// go-clang-dump shows how to dump the AST of a C/C++ file via the Cursor
// visitor API.
//
// ex:
// $ go-clang-dump -fname=foo.cxx
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zimmski/go-clang-phoenix"
)

var fname = flag.String("fname", "", "the file to analyze")

func main() {
	fmt.Printf(":: go-clang-dump...\n")
	flag.Parse()
	fmt.Printf(":: fname: %s\n", *fname)
	fmt.Printf(":: args: %v\n", flag.Args())

	if *fname == "" {
		flag.Usage()
		fmt.Printf("please provide a file name to analyze\n")

		os.Exit(1)
	}

	idx := phoenix.NewIndex(0, 1)
	defer idx.Dispose()

	args := []string{}
	if len(flag.Args()) > 0 && flag.Args()[0] == "-" {
		args = make([]string, len(flag.Args()[1:]))
		copy(args, flag.Args()[1:])
	}

	tu := idx.ParseTranslationUnit(*fname, args, nil, 0)
	defer tu.Dispose()

	fmt.Printf("tu: %s\n", tu.Spelling())
	cursor := tu.TranslationUnitCursor()
	fmt.Printf("cursor-isnull: %v\n", cursor.IsNull())
	fmt.Printf("cursor: %s\n", cursor.Spelling())
	fmt.Printf("cursor-kind: %s\n", cursor.Kind().Spelling())

	fmt.Printf("tu-fname: %s\n", tu.File(*fname).Name())

	fct := func(cursor, parent phoenix.Cursor) phoenix.ChildVisitResult {
		if cursor.IsNull() {
			fmt.Printf("cursor: <none>\n")

			return phoenix.ChildVisit_Continue
		}

		fmt.Printf("%s: %s (%s)\n", cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR())

		switch cursor.Kind() {
		case phoenix.Cursor_ClassDecl, phoenix.Cursor_EnumDecl, phoenix.Cursor_StructDecl, phoenix.Cursor_Namespace:
			return phoenix.ChildVisit_Recurse
		}

		return phoenix.ChildVisit_Continue
	}

	cursor.Visit(fct)

	fmt.Printf(":: bye.\n")
}
