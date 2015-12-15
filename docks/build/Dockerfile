## goclang/build
## a Docker environment to build and test go-clang
from goclang/base
maintainer go-clang-devs@github.com/go-clang

env LLVM_VERSION 3.4

## add repositories
run apt-get install -y software-properties-common && \
	curl -O -L http://llvm.org/apt/llvm-snapshot.gpg.key && \
	apt-key add llvm-snapshot.gpg.key && \
    add-apt-repository --enable-source \
		"deb http://llvm.org/apt/${CODENAME}/ llvm-toolchain-${CODENAME}-${LLVM_VERSION} main" && \
	/bin/rm llvm-snapshot.gpg.key


## update and upgrade
run apt-get update -y && \
	apt-get -V upgrade -y && \
	apt-get -V autoremove -y

## install needed packages
run apt-get -V install -y \
	clang-${LLVM_VERSION} git libclang-${LLVM_VERSION}-dev \
	llvm-${LLVM_VERSION}-tools make \
	gcc g++

## cleanup
run apt-get clean -y

## setup LLVM and CLang (TODO: remove. see go-clang/design#3 go-clang/gen#102)
run ln -s /usr/bin/llvm-config-${LLVM_VERSION} /usr/bin/llvm-config && \
	ln -s /usr/lib/x86_64-linux-gnu/libclang-${LLVM_VERSION}.so \
		  /usr/lib/x86_64-linux-gnu/libclang.so

## setup for go-clang
workdir ${GOPATH}/github.com/go-clang/gen

