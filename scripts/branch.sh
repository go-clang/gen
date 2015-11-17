#!/bin/bash

if [ -z "$1" ]; then
	exit
fi

export LLVM_VERSION=$1

git checkout bootstrap
git checkout -b "v${LLVM_VERSION//.}"

# Generate and install new Clang version
~/go-clang-generate
make install

# Change versions in files
sed -i -e "s/3.4/${LLVM_VERSION}/g" .travis.yml

# Add new generation to git
git add clang-c/
git add *gen*
git add -u # do not forget to use the correct Clang version in TravisCI
