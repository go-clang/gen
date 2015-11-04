package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	// "github.com/termie/go-shutil"
)

type LLVMVersion struct {
	Major    int
	Minor    int
	Subminor int
}

func main() {
	rawLLVMVersion, _, err := execToBuffer("llvm-config", "--version")
	if err != nil {
		exitWithFatal("Cannot determine LLVM version", err)
	}

	matchLLVMVersion := regexp.MustCompile(`^(\d+)\.(\d+)(?:\.(\d+))?`).FindSubmatch(rawLLVMVersion)
	if matchLLVMVersion == nil {
		exitWithFatal("Cannot parse LLVM version", nil)
	}

	var llvmVersion LLVMVersion

	llvmVersion.Major, _ = strconv.Atoi(string(matchLLVMVersion[1]))
	llvmVersion.Minor, _ = strconv.Atoi(string(matchLLVMVersion[2]))
	llvmVersion.Subminor, _ = strconv.Atoi(string(matchLLVMVersion[3]))

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

	clangArguments := []string{
		"-I", ".", // Include current folder
	}

	for _, d := range []string{
		"/usr/local/lib/clang",
		"/usr/include/clang",
	} {
		for _, di := range []string{
			d + fmt.Sprintf("/%d.%d.%d/include", llvmVersion.Major, llvmVersion.Minor, llvmVersion.Subminor),
			d + fmt.Sprintf("/%d.%d/include", llvmVersion.Major, llvmVersion.Minor),
		} {
			if dirExists(di) == nil {
				clangArguments = append(clangArguments, "-I", di)
			}
		}
	}

	fmt.Printf("Using clang arguments: %v\n", clangArguments)

	fmt.Printf("Will generate go-clang for LLVM version %d.%d in current directory\n", llvmVersion.Major, llvmVersion.Minor)

	// TODO reenable
	/*// Copy the Clang-C include directory into the current directory
	_ = os.RemoveAll("./clang-c/")
	if err := shutil.CopyTree(clangCIncludeDir, "./clang-c/", nil); err != nil {
		exitWithFatal(fmt.Sprintf("Cannot copy Clang-C include directory %q into current directory", clangCIncludeDir), err)
	}*/

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

	clangCDirectory := "./clang-c/"
	headers, err := ioutil.ReadDir(clangCDirectory)
	if err != nil {
		exitWithFatal("Cannot list clang-c directory", err)
	}
	for _, h := range headers {
		if h.IsDir() {
			continue
		}

		newHeaderFile(clangCDirectory + h.Name()).handleHeaderFile(clangArguments)
	}
}
