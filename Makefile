.PHONY: all test 

export CC := clang
export CXX := clang++

LLVM_CONFIG ?= llvm-config
CGO_CFLAGS=
CGO_LDFLAGS=$(strip -L$(shell ${LLVM_CONFIG} --libdir) -Wl,-rpath,$(shell ${LLVM_CONFIG} --libdir))

all: test

test:
	CGO_CFLAGS='${CGO_CFLAGS}' CGO_LDFLAGS='${CGO_LDFLAGS}' go test -v -race ./...

coverage:
	CGO_CFLAGS='${CGO_CFLAGS}' CGO_LDFLAGS='${CGO_LDFLAGS}' go test -v -covermode=atomic -coverpkg=./... -coverprofile=coverage.out ./...
