.PHONY: all clean install lint

ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export ROOT_DIR

all: install

clean:
	rm *_gen.go
install:
	CGO_CFLAGS="-I`llvm-config --includedir`" CGO_LDFLAGS="-L`llvm-config --libdir`" go install ./...
lint: install
	errcheck ./... || true
	golint ./... | grep --invert-match -P "(_string.go:)" || true
	go tool vet -all=true -v=true $(ROOT_DIR)/ 2>&1 | grep --invert-match -P "(Checking file|\%p of wrong type|can't check non-constant format)" || true
