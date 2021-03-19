.PHONY: all install install-dependencies install-tools lint test test-full test-verbose

export ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

export CC := clang
export CXX := clang++

ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(ARGS):;@:) # turn arguments into do-nothing targets
export ARGS

all: install-dependencies install-tools install test

install:
	CGO_LDFLAGS="-L`llvm-config --libdir`" go install ./...
install-dependencies:
	go get -u golang.org/x/tools/imports/...
	go get -u github.com/stretchr/testify/...
	go get -u github.com/termie/go-shutil/...

	CGO_LDFLAGS="-L`llvm-config --libdir`" go get github.com/go-clang/bootstrap/...
install-tools:
	# Install linting tools
	go get -u golang.org/x/lint/golint
	go get -u github.com/kisielk/errcheck/...

	# Install code coverage tools
	go get -u golang.org/x/tools/cmd/cover/...
	go get -u github.com/onsi/ginkgo/ginkgo/...
	go get -u github.com/modocache/gover/...
	go get -u github.com/mattn/goveralls/...
lint: install
	$(ROOT_DIR)/scripts/lint.sh
test:
	CGO_LDFLAGS="-L`llvm-config --libdir`" go test -timeout 60s ./...
test-full:
	$(ROOT_DIR)/scripts/test-full.sh
test-verbose:
	CGO_LDFLAGS="-L`llvm-config --libdir`" go test -timeout 60s -v ./...

LLVM_CONFIG ?= llvm-config
# CGO_CFLAGS=$(strip $(shell ${LLVM_CONFIG} --cflags | awk '{$$1=""}1')) -Wno-deprecated-declarations -Wno-language-extension-token -Wno-unused-variable
CGO_CFLAGS=
CGO_LDFLAGS=$(strip -L$(shell ${LLVM_CONFIG} --libdir) -Wl,-rpath,$(shell ${LLVM_CONFIG} --libdir))

pkg/%:
	CGO_CFLAGS='${CGO_CFLAGS}' CGO_LDFLAGS='${CGO_LDFLAGS}' go-install-pkg || true

build/%:
	CGO_CFLAGS='${CGO_CFLAGS}' CGO_LDFLAGS='${CGO_LDFLAGS}' go build -v -x -o go-clang-gen-$* ./cmd/go-clang-gen

gen/%: build/%
	rm -rf $(shell go env GOPATH)/src/github.com/go-clang/clang-v$*/clang; \
		mkdir -p $(shell go env GOPATH)/src/github.com/go-clang/clang-v$*/clang; \
		cd $(shell go env GOPATH)/src/github.com/go-clang/clang-v$*/clang; \
		$(shell go env GOPATH)/src/github.com/go-clang/gen/go-clang-gen-$* -llvm-root=$(shell ${LLVM_CONFIG} --prefix)

docker/build/%:
	docker image build --rm --build-arg LLVM_VERSION=$*.0 -t gcr.io/go-clang/gen:$* ./docker

docker/gen/%: docker/build/%
	docker container run --rm -it -v $(shell go env GOPATH)/src/github.com/go-clang/gen:/go/src/github.com/go-clang/gen -v $(shell go env GOPATH)/src/github.com/go-clang/clang-v$*:/go/src/github.com/go-clang/clang-v$* -w /go/src/github.com/go-clang/gen goclang/gen:$* make gen/$*
