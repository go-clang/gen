package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/template"
)

// File represents a generation file
type File struct {
	Name string

	Imports map[string]struct{}

	Functions []string
}

var templateGenerateFile = template.Must(template.New("go-clang-generate-file").Parse(`package phoenix

// #include "go-clang.h"
import "C"
{{if $.Imports}}
import (
{{range $import, $empty := $.Imports}}	"{{$import}}"
{{end}}){{end}}
{{range $i, $f := $.Functions}}
{{$f}}
{{end}}
`))

func generateFile(f *File) error {
	// TODO remove this hack
	for _, m := range f.Functions {
		if strings.Contains(m, "time.Time") {
			f.Imports["time"] = struct{}{}
		}
		if strings.Contains(m, "unsafe.") {
			f.Imports["unsafe"] = struct{}{}
		}
	}

	var b bytes.Buffer
	if err := templateGenerateFile.Execute(&b, f); err != nil {
		return err
	}

	return ioutil.WriteFile(strings.ToLower(f.Name)+"_gen.go", b.Bytes(), 0600)
}
