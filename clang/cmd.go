package clang

import (
	"bytes"
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-clang/gen"
)

// cmdFatal returns a command error
func cmdFatal(msg string, err error) error {
	if err == nil {
		return fmt.Errorf("FATAL %s", msg)
	} else {
		return fmt.Errorf("FATAL %s: %s", msg, err)
	}
}

//go:embed clang/*
var f embed.FS

const embedDir = "clang"

var fileList = []string{
	"cgoflags.go.txt",
	"cgoflags_dynamic.go.txt",
	"cgoflags_static.go.txt",
	"clang_test.go.txt",
	"completion_test.go.txt",
	"cursor.c.txt",
	"cursor.go.txt",
	"cxstring.go.txt",
	"go-clang.h.txt",
	"translationunit.go.txt",
	"unsavedfile.go.txt",
}

// Cmd executes a generic go-clang-generate command
func Cmd(llvmRoot string, api *gen.API) error {
	llvmConfigPath := filepath.Join(llvmRoot, "bin", "llvm-config")
	if err := fileExists(llvmConfigPath); err != nil {
		return err
	}

	rawLLVMVersion, _, err := execToBuffer(llvmConfigPath, "--version")
	if err != nil {
		return cmdFatal("Cannot determine LLVM version", err)
	}

	llvmVersion := ParseVersion(rawLLVMVersion)
	if llvmVersion == nil {
		return cmdFatal("Cannot parse LLVM version", nil)
	}

	fmt.Println("Found LLVM version", llvmVersion.String())

	rawLLVMIncludeDir, _, err := execToBuffer(llvmConfigPath, "--includedir")
	if err != nil {
		return cmdFatal("Cannot determine LLVM include directory", err)
	}

	clangCIncludeDir := strings.TrimSpace(string(rawLLVMIncludeDir)) + "/clang-c/"
	if err := dirExists(clangCIncludeDir); err != nil {
		return cmdFatal(fmt.Sprintf("Cannot find Clang-C include directory %q", clangCIncludeDir), err)
	}

	fmt.Println("Clang-C include directory", clangCIncludeDir)

	clangArguments := []string{
		"-I", ".", // Include the current directory
	}

	// Find Clang's include directory
	clangResourceDir, _, err := execToBuffer(filepath.Join(llvmRoot, "bin", "clang"), "--print-resource-dir")
	if err == nil { // ignore error for --print-resource-dir flag not yet implements
		clangResourceIncludeDir := filepath.Join(strings.TrimSpace(string(clangResourceDir)), "include")
		if err := dirExists(clangResourceIncludeDir); err != nil {
			return cmdFatal("not fonud clang resource include directory", err)
		}

		clangArguments = append(clangArguments, "-I", clangResourceIncludeDir)
	}

	api.ClangArguments = append(clangArguments, api.ClangArguments...)

	fmt.Printf("Using clang arguments: %v\n", api.ClangArguments)

	fmt.Printf("Will generate go-clang for LLVM version %s into the current directory\n", llvmVersion.String())

	clangCDirectory := "./clang-c/"

	// Copy the Clang-C include directory into the current directory
	_ = os.RemoveAll(clangCDirectory)
	if err := copyTree(clangCIncludeDir, clangCDirectory); err != nil {
		return cmdFatal(fmt.Sprintf("Cannot copy Clang-C include directory %q into current directory", clangCIncludeDir), err)
	}

	clangDirectory := "./"

	// Remove all generated .go files
	if files, err := ioutil.ReadDir(clangDirectory); err != nil {
		return cmdFatal("Cannot read current directory", err)
	} else {
		for _, f := range files {
			fn := f.Name()

			if !f.IsDir() && strings.HasSuffix(fn, "_gen.go") {
				if err := os.Remove(fn); err != nil {
					return cmdFatal(fmt.Sprintf("Cannot remove generated file %q", fn), err)
				}
			}
		}
	}

	// write no generated file into clang directory
	for _, file := range fileList {
		if err := writeEmbedFile(clangDirectory, filepath.Join(embedDir, file)); err != nil {
			return err
		}
	}

	// write clang/doc.go
	clangDoc, err := f.ReadFile(filepath.Join(embedDir, "doc.go.txt"))
	if err != nil {
		return err
	}
	clangDoc = bytes.ReplaceAll(clangDoc, []byte("VERSION"), []byte(strconv.FormatInt(int64(llvmVersion.Major), 10)))
	if err := os.WriteFile(filepath.Join(clangDirectory, "doc.go"), clangDoc, 0644); err != nil {
		return err
	}

	// write clang/clang-c/doc.go
	clangCDoc, err := f.ReadFile(filepath.Join(embedDir, "clang_c.go.txt"))
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(clangCDirectory, "doc.go"), clangCDoc, 0644); err != nil {
		return err
	}

	// write testdata
	dirent, err := f.ReadDir(filepath.Join(embedDir, "testdata"))
	if err != nil {
		return err
	}
	testdataDir := filepath.Join(clangDirectory, "testdata")
	if err := os.Mkdir(testdataDir, 0755); err != nil {
		return err
	}
	for _, ent := range dirent {
		if err := writeEmbedFile(testdataDir, filepath.Join(embedDir, "testdata", ent.Name())); err != nil {
			return err
		}
	}

	headerFiles, err := api.HandleDirectory(clangCDirectory)
	if err != nil {
		return err
	}

	generator := gen.NewGeneration(api)
	generator.AddHeaderFiles(headerFiles)
	if err = generator.Generate(); err != nil {
		return err
	}

	return nil
}

func writeEmbedFile(dir, name string) error {
	data, err := f.ReadFile(name)
	if err != nil {
		return err
	}

	name = strings.TrimSuffix(name, ".txt")
	if err := os.WriteFile(filepath.Join(dir, filepath.Base(name)), data, 0644); err != nil {
		return err
	}

	return nil
}
