#!/bin/bash

set -exuo pipefail

if [ -z "$1" ]; then
	exit
fi

export CODENAME=$(lsb_release --codename --short)
export LLVM_VERSION=$1

# Add repositories
sudo add-apt-repository --enable-source "deb http://llvm.org/apt/${CODENAME}/ llvm-toolchain-${CODENAME}-${LLVM_VERSION} main"
sudo apt-get update

sudo rm /usr/bin/llvm-config
sudo apt-get install -y clang-$LLVM_VERSION libclang1-$LLVM_VERSION libclang-$LLVM_VERSION-dev llvm-$LLVM_VERSION llvm-$LLVM_VERSION-dev llvm-$LLVM_VERSION-runtime libclang-common-$LLVM_VERSION-dev
sudo ln -s /usr/bin/llvm-config-$LLVM_VERSION /usr/bin/llvm-config
