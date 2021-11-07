package clang

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
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

func isSymlink(fi os.FileInfo) bool {
	return (fi.Mode() & os.ModeSymlink) == os.ModeSymlink
}

func copyFile(src, dst string) error {
	srcInfo, _ := os.Stat(src)
	dstInfo, _ := os.Stat(dst)
	if os.SameFile(srcInfo, dstInfo) {
		return fmt.Errorf("%q and %q are the same file", src, dst)
	}

	srcStat, err := os.Lstat(src)
	if err != nil {
		return err
	}

	if _, err := os.Stat(dst); err != nil && !os.IsNotExist(err) {
		return err
	}

	// If we don't follow symlinks and it's a symlink, just link it and be done
	if isSymlink(srcStat) {
		return os.Symlink(src, dst)
	}

	// If we are a symlink, follow it
	if isSymlink(srcStat) {
		src, err = os.Readlink(src)
		if err != nil {
			return err
		}
		srcStat, err = os.Stat(src)
		if err != nil {
			return err
		}
	}

	// Do the actual copy
	fsrc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fsrc.Close()

	fdst, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fdst.Close()

	size, err := io.Copy(fdst, fsrc)
	if err != nil {
		return err
	}

	if size != srcStat.Size() {
		return fmt.Errorf("%s: %d/%d copied", src, size, srcStat.Size())
	}

	return nil
}

func copyMode(src, dst string) error {
	srcStat, err := os.Lstat(src)
	if err != nil {
		return err
	}

	dstStat, err := os.Lstat(dst)
	if err != nil {
		return err
	}

	// They are both symlinks and we can't change mode on symlinks.
	if isSymlink(srcStat) && isSymlink(dstStat) {
		return nil
	}

	// Atleast one is not a symlink, get the actual file stats
	srcStat, _ = os.Stat(src)
	err = os.Chmod(dst, srcStat.Mode())
	return err
}

func copyFunc(src, dst string) (string, error) {
	dstInfo, err := os.Stat(dst)

	if err == nil && dstInfo.Mode().IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}

	if err != nil && !os.IsNotExist(err) {
		return dst, err
	}

	if err := copyFile(src, dst); err != nil {
		return dst, err
	}

	if err := copyMode(src, dst); err != nil {
		return dst, err
	}

	return dst, nil
}

func copyTree(src, dst string) error {
	srcFileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcFileInfo.IsDir() {
		return fmt.Errorf("%q is not a directory", src)
	}

	if _, err := os.Open(dst); !os.IsNotExist(err) {
		return fmt.Errorf("%q already exists", dst)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcFileInfo.Mode()); err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		entryFileInfo, err := os.Lstat(srcPath)
		if err != nil {
			return err
		}

		switch {
		case isSymlink(entryFileInfo):
			linkTo, err := os.Readlink(srcPath)
			if err != nil {
				return err
			}
			// ignore dangling symlink if flag is on
			_, err = os.Stat(linkTo)
			if err != nil {
				return err
			}
			_, err = copyFunc(srcPath, dstPath)
			if err != nil {
				return err
			}
		case entryFileInfo.IsDir():
			err = copyTree(srcPath, dstPath)
			if err != nil {
				return err
			}
		default:
			_, err = copyFunc(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
