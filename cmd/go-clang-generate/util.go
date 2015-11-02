package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unicode"
)

// TODO util.go is just an ugly name...

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

func lowerFirstCharacter(s string) string {
	r := []rune(s)

	r[0] = unicode.ToLower(r[0])

	return string(r)
}

func upperFirstCharacter(s string) string {
	r := []rune(s)

	r[0] = unicode.ToUpper(r[0])

	return string(r)
}
