package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// CmdFatal returns a command error
func CmdFatal(msg string, err error) error {
	if err == nil {
		return fmt.Errorf("FATAL %s", msg)
	} else {
		return fmt.Errorf("FATAL %s: %s", msg, err)
	}
}

// Cmd executes a generic go-clang-generate command
func Cmd(args []string) error {
	rawLLVMVersion, _, err := execToBuffer("llvm-config", "--version")
	if err != nil {
		return CmdFatal("Cannot determine LLVM version", err)
	}

	llvmVersion := ParseVersion(rawLLVMVersion)
	if llvmVersion == nil {
		return CmdFatal("Cannot parse LLVM version", nil)
	}

	fmt.Println("Found LLVM version", llvmVersion.String())

	rawLLVMIncludeDir, _, err := execToBuffer("llvm-config", "--includedir")
	if err != nil {
		return CmdFatal("Cannot determine LLVM include directory", err)
	}

	clangCIncludeDir := strings.TrimSpace(string(rawLLVMIncludeDir)) + "/clang-c/"
	if err := dirExists(clangCIncludeDir); err != nil {
		return CmdFatal(fmt.Sprintf("Cannot find Clang-C include directory %q", clangCIncludeDir), err)
	}

	fmt.Println("Clang-C include directory", clangCIncludeDir)

	clangArguments := []string{
		"-I", ".", // Include the current directory
	}

	// Find Clang's include directory
	for _, d := range []string{
		"/usr/local/lib/clang",
		"/usr/include/clang",
	} {
		for _, di := range []string{
			d + fmt.Sprintf("/%s/include", llvmVersion.String()),
			d + fmt.Sprintf("/%s/include", llvmVersion.StringMinor()),
		} {
			if dirExists(di) == nil {
				clangArguments = append(clangArguments, "-I", di)
			}
		}
	}

	fmt.Printf("Using clang arguments: %v\n", clangArguments)

	fmt.Printf("Will generate go-clang for LLVM version %s into the current directory\n", llvmVersion.String())

	clangCDirectory := "./clang-c/"
	// TODO reenable
	/*// Copy the Clang-C include directory into the current directory
	_ = os.RemoveAll(clangCDirectory)
	if err := shutil.CopyTree(clangCIncludeDir, clangCDirectory, nil); err != nil {
		return CmdFatal(fmt.Sprintf("Cannot copy Clang-C include directory %q into current directory", clangCIncludeDir), err)
	}*/

	// Remove all generated .go files
	if files, err := ioutil.ReadDir("./"); err != nil {
		return CmdFatal("Cannot read current directory", err)
	} else {
		for _, f := range files {
			fn := f.Name()

			if !f.IsDir() && strings.HasSuffix(fn, "_gen.go") {
				if err := os.Remove(fn); err != nil {
					return CmdFatal(fmt.Sprintf("Cannot remove generated file %q", fn), err)
				}
			}
		}
	}

	headers, err := ioutil.ReadDir(clangCDirectory)
	if err != nil {
		return CmdFatal("Cannot list clang-c directory", err)
	}
	for _, h := range headers {
		if h.IsDir() && !strings.HasSuffix(h.Name(), ".h") {
			continue
		}

		if err := HandleHeaderFile(clangCDirectory+h.Name(), clangArguments); err != nil {
			return err
		}
	}

	return nil
}
