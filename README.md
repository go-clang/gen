# go-clang-phoenix [![GoDoc](https://godoc.org/github.com/zimmski/go-clang-phoenix?status.png)](https://godoc.org/github.com/zimmski/go-clang-phoenix) [![Build Status](https://travis-ci.org/zimmski/go-clang-phoenix.svg?branch=master)](https://travis-ci.org/zimmski/go-clang-phoenix) [![Coverage Status](https://coveralls.io/repos/zimmski/go-clang-phoenix/badge.png?branch=master)](https://coveralls.io/r/zimmski/go-clang-phoenix?branch=master)

Native Go bindings for the C API of clang.

## Installation

```bash
CGO_CFLAGS="-I`llvm-config --includedir`" \
  CGO_LDFLAGS="-L`llvm-config --libdir`" \
  go get github.com/zimmski/go-clang-phoenix
```

## Example

An example on how to use the AST visior of clang can be found in [/cmd/go-clang-dump/main.go](/cmd/go-clang-dump/main.go)
