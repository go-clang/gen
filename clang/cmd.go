package clang

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-clang/gen"
)

//go:embed _embed
var embedClang embed.FS

const (
	embedDirRootPath = "_embed"

	clangDirName    = "clang"
	clangCDirName   = "clang-c"
	testdataDirName = "testdata"
)

var (
	embedClangDirPath    = filepath.Join(embedDirRootPath, clangDirName)
	embedClanCgDirPath   = filepath.Join(embedDirRootPath, clangDirName, clangCDirName)
	embedTestdataDirPath = filepath.Join(embedDirRootPath, testdataDirName)
)

// Cmd executes a generic go-clang-generate command.
func Cmd(llvmRoot string, api *gen.API) error {
	llvmConfigPath := filepath.Join(llvmRoot, "bin", "llvm-config")
	if err := fileExists(llvmConfigPath); err != nil {
		return err
	}

	rawLLVMVersion, _, err := execToBuffer(llvmConfigPath, "--version")
	if err != nil {
		return fmt.Errorf("cannot determine LLVM version: %w", err)
	}

	llvmVersion := ParseVersion(rawLLVMVersion)
	if llvmVersion == nil {
		return errors.New("cannot parse LLVM version")
	}
	fmt.Printf("detected the LLVM version: %s\n", llvmVersion)

	rawLLVMIncludeDir, _, err := execToBuffer(llvmConfigPath, "--includedir")
	if err != nil {
		return fmt.Errorf("cannot determine LLVM include directory: %w", err)
	}

	clangCIncludeDir := filepath.Join(strings.TrimSpace(string(rawLLVMIncludeDir)), clangCDirName)
	if err := dirExists(clangCIncludeDir); err != nil {
		return fmt.Errorf("cannot find %q include directory: %w", clangCIncludeDir, err)
	}
	fmt.Printf("found clang-c include directory: %s\n", clangCIncludeDir)

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get working directory: %w", err)
	}
	clangDirPath := filepath.Join(cwd, clangDirName)
	clangCDirPath := filepath.Join(cwd, clangDirName, clangCDirName)
	testdataDirPath := filepath.Join(cwd, testdataDirName)

	clangArguments := []string{
		"-I" + clangDirPath, // include clang directory
	}

	// find Clang's resource directory
	clangResourceDir, _, err := execToBuffer(filepath.Join(llvmRoot, "bin", "clang"), "--print-resource-dir")
	if err == nil { // Clang have --print-resource-dir flag
		clangResourceIncludeDir := filepath.Join(strings.TrimSpace(string(clangResourceDir)), "include")
		if err := dirExists(clangResourceIncludeDir); err != nil {
			return fmt.Errorf("not fonud clang resource directory: %w", err)
		}
		clangArguments = append(clangArguments, "-I"+clangResourceIncludeDir)
	} else {
		llvmLibDir, _, err := execToBuffer(llvmConfigPath, "--libdir")
		if err != nil {
			return fmt.Errorf("cannot determine LLVM libdir directory: %w", err)
		}
		clangArguments = append(clangArguments, "-I"+filepath.Join(strings.TrimSpace(string(llvmLibDir)), "clang", llvmVersion.String(), "include"))
	}

	// set ClangArguments
	api.ClangArguments = append(api.ClangArguments, clangArguments...)

	fmt.Printf("using clang arguments: %v\n", api.ClangArguments)
	fmt.Printf("will generate go-clang for %s version into the ./%s directory\n", llvmVersion, clangDirName)

	// remove all generated _gen.go files
	oldGenFiles, err := os.ReadDir(clangDirPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cannot read %s directory: %w", clangDirName, err)
	}
	for _, f := range oldGenFiles {
		fname := f.Name()
		if !f.IsDir() && strings.HasSuffix(fname, "_gen.go") {
			if err := os.Remove(fname); err != nil {
				return fmt.Errorf("cannot remove %q generated file: %w", fname, err)
			}
		}
	}
	_ = os.RemoveAll(clangCDirPath)   // remove old clang/clang-c directory
	_ = os.RemoveAll(testdataDirPath) // remove old testdata directory

	// copy the clang-c include directory into the clang directory
	if err := copyTree(clangCIncludeDir, clangCDirPath); err != nil {
		return fmt.Errorf("cannot copy Clang C bindings %q include headers directory into %s directory: %w", clangCIncludeDir, clangDirName, err)
	}

	// write non-generated file into clang directory
	if err := WriteEmbedFile(clangDirPath, embedClangDirPath); err != nil {
		return fmt.Errorf("could not write embedded %s non-generated file: %w", embedClangDirPath, err)
	}

	// write testdata files
	if err := WriteEmbedFile(testdataDirPath, embedTestdataDirPath); err != nil {
		return fmt.Errorf("could not write embedded %s testdata file: %w", embedTestdataDirPath, err)
	}

	// analyze LLVM version before replacing the import path for clang/doc.go to support the LLVM 3.x family release policy
	// in the LLVM 3.x era, the release version combined the major version and the minor version
	var replaceLLVMVersion string
	switch llvmVersion.Major {
	case 3: // 3.4.0 ~ 3.9.0
		replaceLLVMVersion = llvmVersion.StringMinor()
	default:
		replaceLLVMVersion = strconv.Itoa(llvmVersion.Major)
	}

	// write clang/doc.go
	docData := strings.ReplaceAll(clangDocTmpl, replaceMark, replaceLLVMVersion)
	clangDocPath := filepath.Join(clangDirPath, "doc.go")
	if err := os.WriteFile(clangDocPath, []byte(docData), 0o644); err != nil {
		return fmt.Errorf("could not write %s file: %w", clangDocPath, err)
	}

	// write clang/clang-c/doc.go
	clangCDocPath := filepath.Join(clangCDirPath, "doc.go")
	if err := os.WriteFile(clangCDocPath, []byte(clangCDocTmpl), 0o644); err != nil {
		return fmt.Errorf("could not write %s file: %w", clangCDocPath, err)
	}

	// handle Clang headers
	headerFiles, err := api.HandleDirectory("./" + clangDirName + string(os.PathSeparator) + clangCDirName)
	if err != nil {
		return fmt.Errorf("could not handle clang-c header directory: %w", err)
	}

	// initialize generator
	generator := gen.NewGeneration(api)
	generator.AddHeaderFiles(headerFiles)

	// generation Clang binding
	if err := generator.Generate(); err != nil {
		return fmt.Errorf("could not generate: %w", err)
	}

	return nil
}

const replaceMark = "$VERSION$"

const clangDocTmpl = `// Package clang provides the Clang C API bindings for Go.
package clang

import (
	_ "github.com/go-clang/clang-v` + replaceMark + `/clang/clang-c"
)
`

const clangCDocTmpl = `// Package clang_c holds clang binding C header files.
package clang_c
`

// WriteEmbedFile reads embedPath file from embedClang and writes to dstDir.
//
// The embedPath must be a full path from the embed root directory.
func WriteEmbedFile(dstDir, embedDir string) error {
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return fmt.Errorf("make %s directory: %w", dstDir, err)
	}

	d, err := embedClang.ReadDir(embedDir)
	if err != nil {
		return fmt.Errorf("unebla to read %s directory: %w", embedDir, err)
	}

	for _, ent := range d {
		if ent.IsDir() {
			// not handle directory, this function only write file
			continue
		}

		fname := ent.Name()
		data, err := embedClang.ReadFile(filepath.Join(embedDir, fname))
		if err != nil {
			return fmt.Errorf("could not read embedded %s file: %w", fname, err)
		}

		if err := os.WriteFile(filepath.Join(dstDir, fname), data, 0o644); err != nil {
			return fmt.Errorf("could not write %s file: %w", fname, err)
		}
	}

	return nil
}
