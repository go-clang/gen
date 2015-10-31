.PHONY: all clean generate install lint test test-verbose

ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export ROOT_DIR

all: install test

clean:
	rm -r clang-c/
	rm *_gen.go
generate:
	CGO_CFLAGS="-I`llvm-config --includedir`" CGO_LDFLAGS="-L`llvm-config --libdir`" go run cmd/go-clang-generate/*
install:
	CGO_CFLAGS="-I`llvm-config --includedir`" CGO_LDFLAGS="-L`llvm-config --libdir`" go install ./...
lint: install
	scripts/lint.sh
test:
	CGO_CFLAGS="-I`llvm-config --includedir`" CGO_LDFLAGS="-L`llvm-config --libdir`" go test -timeout 60s -race ./...
test-verbose:
	CGO_CFLAGS="-I`llvm-config --includedir`" CGO_LDFLAGS="-L`llvm-config --libdir`" go test -timeout 60s -race -v ./...
