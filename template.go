// Code generated by go generate; DO NOT EDIT.
package main

import (
	"os"
	"text/template"
	"path/filepath"
)

type dirMetaData struct {
	path     string
	fileMode os.FileMode
}

type fileMetaData struct {
	path     string
	fileMode os.FileMode
	content  string
}

var (
	templateDirs = []dirMetaData{
	}

	templateFiles = []fileMetaData{
		fileMetaData{
			path: 		".gitignore",
			fileMode: 	os.FileMode(420),
			content:	`server
test.bin
cmd/server/ver.go
`,
		},
		fileMetaData{
			path: 		"main.go",
			fileMode: 	os.FileMode(420),
			content:	`package main

import "fmt"

func main() {
	fmt.Println("{{.GreetMsg}}")
}
`,
		},
	}
)

func GenScaffold(outdir string, data interface{}) error {
	// ensure outputdir
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		if err := os.Mkdir(outdir, os.ModePerm); err != nil {
			return err
		}
	}

	// prepare directories
	for _, di := range templateDirs {
		dir := filepath.Join(outdir, di.path)
		if err := os.MkdirAll(dir, di.fileMode); err != nil {
			return err
		}
	}

	// prepare files
	for _, fi := range templateFiles {
		path := filepath.Join(outdir, fi.path)
		f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fi.fileMode)
		if err != nil {
			return err
		}
		defer f.Close()

		t, err := template.New("").Parse(fi.content)
		if err != nil {
			return err
		}
		if err := t.Execute(f, data); err != nil {
			return err
		}
	}

	return nil
}