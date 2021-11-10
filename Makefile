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
CGO_CFLAGS=
CGO_LDFLAGS=$(strip -L$(shell ${LLVM_CONFIG} --libdir) -Wl,-rpath,$(shell ${LLVM_CONFIG} --libdir))

go-clang-gen-$*: build/%
build/%:
	CGO_CFLAGS='${CGO_CFLAGS}' CGO_LDFLAGS='${CGO_LDFLAGS}' go build -v -x -o go-clang-gen-$* ./cmd/go-clang-gen

gen/%:
	cd $(shell go env GOPATH)/src/github.com/go-clang/clang-v$*; \
		$(shell go env GOPATH)/src/github.com/go-clang/gen/go-clang-gen-$* -llvm-root=$(shell ${LLVM_CONFIG} --prefix)

GOLANG_VERSION=1.17.2
TARGET=llvm

.PHONY: docker/build/4 docker/build/5 docker/build/6 docker/build/7 docker/build/8 docker/build/9 docker/build/10
docker/build/4 docker/build/5 docker/build/6: UBUNTU_VERSION=18.04
docker/build/7 docker/build/8 docker/build/9 docker/build/10: UBUNTU_VERSION=20.04
docker/build/4 docker/build/5 docker/build/6 docker/build/7 docker/build/8 docker/build/9 docker/build/10:
	docker image build --rm --target=${TARGET} --build-arg UBUNTU_VERSION=${UBUNTU_VERSION} --build-arg LLVM_VERSION=${@F} --build-arg GOLANG_VERSION=${GOLANG_VERSION} -t goclang/base:${@F} -f ./hack/dockerfiles/llvm-4-10.dockerfile .

.PHONY: docker/build/11 docker/build/12 docker/build/13
docker/build/11 docker/build/12 docker/build/13:
	docker image build --rm --target=${TARGET} --build-arg LLVM_VERSION=${@F} --build-arg GOLANG_VERSION=${GOLANG_VERSION} -t goclang/base:${@F} -f ./hack/dockerfiles/llvm-11-13.dockerfile .

docker/gen/%: TARGET=gen
docker/gen/%: docker/build/%
	docker container run --rm -it --mount type=bind,src=$(shell go env GOPATH)/src/github.com/go-clang/clang-v$*,dst=/go/src/github.com/go-clang/clang-v$* -w /go/src/github.com/go-clang/gen goclang/base:$* make gen/$*
