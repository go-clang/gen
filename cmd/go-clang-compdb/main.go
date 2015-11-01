// go-clang-compdb dumps the content of a clang compilation database
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zimmski/go-clang-phoenix"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("**error: you need to give a directory containing a 'compile_commands.json' file\n")

		os.Exit(1)
	}

	dir := os.ExpandEnv(os.Args[1])
	fmt.Printf(":: inspecting [%s]...\n", dir)

	fname := filepath.Join(dir, "compile_commands.json")
	f, err := os.Open(fname)
	if err != nil {
		fmt.Printf("**error: could not open file [%s]: %v\n", fname, err)

		os.Exit(1)
	}
	f.Close()

	err, db := phoenix.CompilationDatabase_fromDirectory(dir)
	if err != nil {
		fmt.Printf("**error: could not open compilation database at [%s]: %v\n", dir, err)

		os.Exit(1)
	}
	defer db.Dispose()

	cmds := db.AllCompileCommands()
	ncmds := cmds.Size()

	fmt.Printf(":: got %d compile commands\n", ncmds)

	for i := uint16(0); i < ncmds; i++ {
		cmd := cmds.Command(i)

		fmt.Printf("::  --- cmd=%d ---\n", i)
		fmt.Printf("::  dir= %q\n", cmd.Directory())

		nargs := cmd.NumArgs()
		fmt.Printf("::  nargs= %d\n", nargs)

		sargs := make([]string, 0, nargs)
		for iarg := uint16(0); iarg < nargs; iarg++ {
			arg := cmd.Arg(iarg)
			sfmt := "%q, "
			if iarg+1 == nargs {
				sfmt = "%q"
			}
			sargs = append(sargs, fmt.Sprintf(sfmt, arg))

		}

		fmt.Printf("::  args= {%s}\n", strings.Join(sargs, ""))
		if i+1 != ncmds {
			fmt.Printf("::\n")
		}
	}
	fmt.Printf(":: inspecting [%s]... [done]\n", dir)
}
