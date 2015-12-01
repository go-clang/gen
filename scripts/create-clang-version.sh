#!/bin/bash

if [ -z "$1" ]; then
	exit
fi

export LLVM_VERSION=$1

# Switch Clang version
$(dirname "$0")/switch-clang-version.sh $LLVM_VERSION || exit

# Initialize the repository
mkdir v${LLVM_VERSION}
cd v${LLVM_VERSION} || exit

git clone https://github.com/go-clang/bootstrap.git . || exit
git remote rename origin bootstrap || exit
git remote add origin git@github.com:go-clang/v${LLVM_VERSION}.git || exit

# Generate the new Clang version
$(dirname "$0")/generate-and-test.sh $LLVM_VERSION || exit
