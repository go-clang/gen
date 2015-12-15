.PHONY: all install install-dependencies install-tools lint test test-all test-verbose

ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export ROOT_DIR

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
install-tools:
	# Install linting tools
	go get -u golang.org/x/tools/cmd/vet/...
	go get -u github.com/golang/lint/...
	go get -u github.com/kisielk/errcheck/...

	# Install code coverage tools
	go get -u golang.org/x/tools/cmd/cover/...
	go get -u github.com/onsi/ginkgo/ginkgo/...
	go get -u github.com/modocache/gover/...
	go get -u github.com/mattn/goveralls/...
lint: install
	scripts/lint.sh
test:
	CGO_LDFLAGS="-L`llvm-config --libdir`" go test -timeout 60s -race ./...
test-verbose:
	CGO_LDFLAGS="-L`llvm-config --libdir`" go test -timeout 60s -race -v ./...

docker-images: docks/build/Dockerfile docks/base/Dockerfile
	docker build --rm -f ./docks/base/Dockerfile  --tag=go-clang/base .
	docker build --rm -f ./docks/build/Dockerfile --tag=go-clang/build .

test-all:
	docker run -v $(shell pwd)/..:/go/src/github.com/go-clang \
		-w /go/src/github.com/go-clang/gen \
		--rm \
		go-clang/build make all lint

