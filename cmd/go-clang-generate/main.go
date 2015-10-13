package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/termie/go-shutil"
)

func execToBuffer(cmd ...string) (out []byte, exitStatus int, err error) {
	c := exec.Command(cmd[0], cmd[1:]...)

	out, err = c.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return out, status.ExitStatus(), err
			}
		}

		return out, 0, err
	}

	return out, 0, nil
}

func exitWithFatal(msg string, err error) {
	if err == nil {
		fmt.Printf("FATAL, %s\n", msg)
	} else {
		fmt.Printf("FATAL, %s: %s\n", msg, err)
	}

	os.Exit(1)
}

func stat(filepath string) (os.FileInfo, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return fi, nil
}

var (
	errNotADirectory = errors.New("not a directory")
	errNotAFile      = errors.New("not a file")
)

func dirExists(path string) error {
	fi, err := stat(path)
	if err != nil {
		return err
	}

	if !fi.Mode().IsDir() {
		return errNotADirectory
	}

	return nil
}

func fileExists(filepath string) error {
	fi, err := stat(filepath)
	if err != nil {
		return err
	}

	if fi.Mode().IsDir() {
		return errNotAFile
	}

	return nil
}

func main() {
	rawLLVMVersion, _, err := execToBuffer("llvm-config", "--version")
	if err != nil {
		exitWithFatal("Cannot determine LLVM version", err)
	}

	matchLLVMVersion := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)`).FindSubmatch(rawLLVMVersion)
	if matchLLVMVersion == nil {
		exitWithFatal("Cannot parse LLVM version", nil)
	}

	var llvmVersion struct {
		Major int
		Minor int
		Patch int
	}

	llvmVersion.Major, _ = strconv.Atoi(string(matchLLVMVersion[1]))
	llvmVersion.Minor, _ = strconv.Atoi(string(matchLLVMVersion[2]))
	llvmVersion.Patch, _ = strconv.Atoi(string(matchLLVMVersion[3]))

	fmt.Println("Found LLVM version", string(matchLLVMVersion[0]))

	rawLLVMIncludeDir, _, err := execToBuffer("llvm-config", "--includedir")
	if err != nil {
		exitWithFatal("Cannot determine LLVM include directory", err)
	}

	clangCIncludeDir := strings.TrimSpace(string(rawLLVMIncludeDir)) + "/clang-c/"
	if err := dirExists(clangCIncludeDir); err != nil {
		exitWithFatal(fmt.Sprintf("Cannot find Clang-C include directory %q", clangCIncludeDir), err)
	}

	fmt.Println("Clang-C include directory", clangCIncludeDir)

	fmt.Printf("Will generate go-clang for LLVM version %d.%d in current directory\n", llvmVersion.Major, llvmVersion.Minor)

	// Copy the Clang-C include directory into the current directory
	_ = os.RemoveAll("clang-c/")
	if err := shutil.CopyTree(clangCIncludeDir, "clang-c/", nil); err != nil {
		exitWithFatal(fmt.Sprintf("Cannot copy Clang-C include directory %q into current directory", clangCIncludeDir), err)
	}

	// Remove all generated .go files
	if files, err := ioutil.ReadDir("./"); err != nil {
		exitWithFatal("Cannot read current directory", err)
	} else {
		for _, f := range files {
			fn := f.Name()

			if !f.IsDir() && strings.HasSuffix(fn, "_gen.go") {
				if err := os.Remove(fn); err != nil {
					exitWithFatal(fmt.Sprintf("Cannot remove generated file %q", fn), err)
				}
			}
		}
	}
}
