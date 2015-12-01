#!/bin/bash

export GO_VERSION=1.4.3

# Install Go
mkdir -p $HOME/go

echo "Downloading Go:"
wget -nv https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz
tar -xf go${GO_VERSION}.linux-amd64.tar.gz -C $HOME/go
rm go${GO_VERSION}.linux-amd64.tar.gz

# Setup user
echo 'export GOPATH=$HOME/go' >> $HOME/.bashrc
echo 'export GOROOT=$GOPATH/go' >> $HOME/.bashrc
echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> $HOME/.bashrc
echo 'cd $HOME/go/src/github.com/go-clang/gen/' >> $HOME/.bashrc

# TODO how can we load .bashrc at this point?
export GOPATH=$HOME/go
export GOROOT=$GOPATH/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
cd $HOME/go/src/github.com/go-clang/gen/

# Install go-clang
make install-dependencies
make install-tools
make install
