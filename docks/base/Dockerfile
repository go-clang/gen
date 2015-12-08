## go-clang/base
## an ubuntu-based docker container with go.
from ubuntu:trusty
maintainer go-clang-devs@github.com/go-clang

env CODENAME trusty
env GO_VERSION 1.4.3

run apt-get update -y && apt-get install -y curl

## install golang
run mkdir -p /go/src
env GOPATH /go
env GOROOT /usr/lib/go
env PATH $PATH:$GOPATH/bin:$GOROOT/bin
run cd /usr/lib && \
	curl -O -L https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
	tar zxf go${GO_VERSION}.linux-amd64.tar.gz && \
	/bin/rm go${GO_VERSION}.linux-amd64.tar.gz && \
	go version

## setup for go-clang
workdir ${GOPATH}/github.com/go-clang/gen


