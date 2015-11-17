#!/bin/bash

if [ -z "$1" ]; then
	exit
fi

export LLVM_VERSION=$1

sudo rm /usr/bin/llvm-config
sudo rm /usr/lib/x86_64-linux-gnu/libclang.so
sudo apt-get install -y clang-$LLVM_VERSION libclang1-$LLVM_VERSION libclang-$LLVM_VERSION-dev llvm-$LLVM_VERSION llvm-$LLVM_VERSION-dev llvm-$LLVM_VERSION-runtime libclang-common-$LLVM_VERSION-dev
sudo ln -s /usr/bin/llvm-config-$LLVM_VERSION /usr/bin/llvm-config
sudo ln -s /usr/lib/x86_64-linux-gnu/libclang-$LLVM_VERSION.so /usr/lib/x86_64-linux-gnu/libclang.so
