#!/bin/bash

if [ -z "$1" ]; then
	exit
fi

export LLVM_VERSION=$1

# Switch Clang version
$(dirname "$0")/switch-clang-version.sh $LLVM_VERSION || exit

# Update the repository
cd go-clang-phoenix-v${LLVM_VERSION} || exit

git checkout bootstrap/master || exit
LAST_BOOTSTRAP=$(git rev-parse HEAD)
git checkout master
git reset --hard $LAST_BOOTSTRAP

git fetch --prune bootstrap || exit
git rebase bootstrap/master master

# Generate the new Clang version
$(dirname "$0")/generate-and-test.sh $LLVM_VERSION || exit
