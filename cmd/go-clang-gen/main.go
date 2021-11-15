package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/go-clang/gen"
	genclang "github.com/go-clang/gen/clang"
	"github.com/go-clang/gen/cmd/go-clang-gen/runtime"
)

var (
	flagLLVMRoot string
)

func init() {
	flag.StringVar(&flagLLVMRoot, "llvm-root", "", "path of llvm root directory")
}

func main() {
	flag.Parse()

	api := &gen.API{
		PrepareFunctionName:     runtime.PrepareFunctionName,
		PrepareFunction:         runtime.PrepareFunction,
		FilterFunction:          runtime.FilterFunction,
		FilterFunctionParameter: runtime.FilterFunctionParameter,
		FixedFunctionName:       runtime.FixedFunctionName,
		PrepareStructFields:     runtime.PrepareStructMembers,
		FilterStructFieldGetter: runtime.FilterStructMemberGetter,
	}

	if flagLLVMRoot == "" {
		c := exec.Command("llvm-config", "--prefix")
		prefix, err := c.CombinedOutput()
		if err != nil {
			if exitErr := new(exec.ExitError); errors.As(err, &exitErr) {
				err = exitErr
			}
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		prefixDir := strings.TrimSpace(string(prefix))
		if rootDir, err := os.Stat(prefixDir); err == nil && rootDir.IsDir() {
			flagLLVMRoot = prefixDir
		}

		if flagLLVMRoot == "" {
			fmt.Fprintln(os.Stderr, "couldn't parse LLVM root directory")
			os.Exit(1)
		}
	}

	if err := genclang.Cmd(flagLLVMRoot, api); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
