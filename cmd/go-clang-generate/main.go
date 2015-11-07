package main

import (
	"fmt"
	"os"

	"github.com/zimmski/go-clang-phoenix/cmd/go-clang-generate/generate"
	generateclang "github.com/zimmski/go-clang-phoenix/cmd/go-clang-generate/generate/clang"
)

func main() {
	api := &generate.API{}

	err := generateclang.Cmd(os.Args[1:], api)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}
