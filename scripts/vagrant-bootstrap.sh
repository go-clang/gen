#!/bin/bash

set -exuo pipefail

export CODENAME=$(lsb_release --codename --short)
export LLVM_VERSION=3.4

# Add repositories
add-apt-repository --enable-source "deb http://llvm.org/apt/${CODENAME}/ llvm-toolchain-${CODENAME}-${LLVM_VERSION} main"
apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 15CF4D18AF4F7421
(wget -O - http://llvm.org/apt/llvm-snapshot.gpg.key | apt-key add -) &> /dev/null

# Update and upgrade
apt-get update
apt-get -V upgrade -y
apt-get -V autoremove -y

# Install needed packages
apt-get -V install -y clang-${LLVM_VERSION} git libclang-${LLVM_VERSION}-dev llvm-${LLVM_VERSION}-tools make

# Setup LLVM and Clang
ln -s /usr/bin/llvm-config-$LLVM_VERSION /usr/bin/llvm-config
ln -s /usr/lib/x86_64-linux-gnu/libclang-$LLVM_VERSION.so /usr/lib/x86_64-linux-gnu/libclang.so

# Change the owner of the go directory. This is needed because the missing parent folders up to our synced folder are created with the user "root" instead of "vagrant".
chown -R vagrant:vagrant /home/vagrant/go
