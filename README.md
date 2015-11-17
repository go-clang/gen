# go-clang-phoenix [![GoDoc](https://godoc.org/github.com/zimmski/go-clang-phoenix?status.png)](https://godoc.org/github.com/zimmski/go-clang-phoenix) [![Build Status](https://travis-ci.org/zimmski/go-clang-phoenix.svg?branch=master)](https://travis-ci.org/zimmski/go-clang-phoenix) [![Coverage Status](https://coveralls.io/repos/zimmski/go-clang-phoenix/badge.png?branch=master)](https://coveralls.io/r/zimmski/go-clang-phoenix?branch=master)

Native Go bindings for the C API of Clang.

## Installation

```bash
CGO_CFLAGS="-I`llvm-config --includedir`" \
  CGO_LDFLAGS="-L`llvm-config --libdir`" \
  go get github.com/zimmski/go-clang-phoenix
```

## Example

An example on how to use the AST visior of Clang can be found in [/cmd/go-clang-dump/main.go](/cmd/go-clang-dump/main.go)

## How to develop for go-clang-phoenix?

You want to contribute to go-clang-phoenix? GREAT! If you are here because of a bug you want to fix or a feature you want to add you can just read on, otherwise we have a list of [open issues in the tracker](/issues). Just choose something you think you can work on and discuss your plans in the issue by commenting on it.

The development branch of go-clang-phoenix is not `master` it is `bootstrap`. `master` just holds the C API bindings of Clang for the latest stable Clang version. We therefore only accept changes based on the `bootstrap` branch.

To ease the development process we have our own development environment based on [Vagrant](https://www.vagrantup.com/). The provided Vagrantfile executed in the root of the repository will setup an Ubuntu VM with our currently used Go version as well as Clang 3.4 and will set up everything that is needed to development and handle new versions of Clang.

### Generate bindings for the current Clang version (VM)

The following command will recompile `go-clang-generate` and will regenerate the bindings for the currently set up Clang version.

```bash
make generate
```

### Switch to a different Clang version (VM)

Replace `3.4` with the Clang version you want to switch to.

```bash
make switch-clang-version 3.4
```

This command will install and configure everything that is needed to develop with the given Clang version. The command will however not generate new bindings for the version.

### Do a PR

Every PR must be prepared using the following commands:

```bash
make switch-clang-version 3.4
make generate
make install
make test
make lint
```

This will generate the bindings with the correct `bootstrap` Clang version, make sure that the bindings compile, run all tests and process the source code with the project's linters. Make sure that you do not introduce new linting problems.

## Maintainer documentation

The following sections are specific to the maintainer process.

### Branch a new Clang version

Every now and then a new Clang version emerges which needs to be generated for go-clang-phoenix. This can be done inside the development VM using the following statement. Replace `3.4` with the Clang version you want to branch off.

```bash
make branch 3.4
```

This will install, configure, generate and install the given Clang version in a branch "v34". Please note, that the "dots" of the version will not be included in the branch name. This is needed to trick gopkg.in believing that this is a new major version. The commit and push for the new version has do be done by hand. The branch should pass all TravisCI checks.
