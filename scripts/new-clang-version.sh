#!/bin/bash

if [ -z "$1" ]; then
	exit
fi

export LLVM_VERSION=$1

# Switch Clang version
$(dirname "$0")/switch-clang-version.sh $LLVM_VERSION || exit

# Initialize the repository
mkdir go-clang-phoenix-v${LLVM_VERSION}
cd go-clang-phoenix-v${LLVM_VERSION} || exit

git clone https://github.com/zimmski/go-clang-phoenix-bootstrap.git . || exit
git remote rename origin bootstrap || exit
git remote add origin git@github.com:zimmski/go-clang-phoenix-v${LLVM_VERSION}.git || exit

# Generate and install new Clang version
cd clang/ || exit

rm -rf clang-c/
rm *_gen.go

go-clang-gen || exit

cd ..

make install || exit
make test || exit

# Change versions in files
sed -i -e "s/3.4/${LLVM_VERSION}/g" .travis.yml

# Show the current state of the repository
git status
