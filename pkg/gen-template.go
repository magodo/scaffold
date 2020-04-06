// +build ignore

// This program generate a go file (template.go) base on the *template* folder, used to generate scaffold from template.
// The output go file path is at the same directory as the invoking go file (which contains the "//go generate" directives).
//
// It accepts following arguments:
// - tempalte dir
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

var outputTemplate string

type dirMetaData struct {
	Path     string
	FileMode os.FileMode
}

type fileMetaData struct {
	Path     string
	Content  string
	FileMode os.FileMode
}

type TemplateEntries struct {
	Dirs  []dirMetaData
	Files []fileMetaData
}

type Engine struct {
	// template root directory
	rootDir string
	// tempalte.go's package name
	InvokerPackage string
	// template content
	TemplateEntries
}

func main() {
	const outputTemplateGoPath = "template.go"
	templateRootDir := os.Args[1]

	// ensure templateRootDir is available
	if _, err := os.Stat(templateRootDir); os.IsNotExist(err) {
		log.Fatal(err)
	}

	invokerPackage := os.Getenv("GOPACKAGE")

	engine := newEngine(templateRootDir, invokerPackage)
	if err := filepath.Walk(templateRootDir, engine.visit); err != nil {
		err = errors.Wrap(err, "failed to walk")
		log.Fatal(err)
	}

	// generate template go file in current directory
	f, err := os.Create(outputTemplateGoPath)
	if err != nil {
		err = errors.Wrapf(err, "failed to create %s", outputTemplateGoPath)
		log.Fatal(err)
	}
	defer f.Close()
	template.Must(template.New("").Parse(outputTemplate)).Execute(f, engine)
}

func newEngine(root string, invokerPkg string) *Engine {
	return &Engine{
		rootDir:         root,
		InvokerPackage:  invokerPkg,
		TemplateEntries: TemplateEntries{},
	}
}

func (e *Engine) visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	return e.addTemplateEntry(path, info)
}

func (e *Engine) addTemplateEntry(path string, info os.FileInfo) error {
	templateRelPath, err := filepath.Rel(e.rootDir, path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		if templateRelPath != "." {
			e.TemplateEntries.Dirs = append(e.TemplateEntries.Dirs,
				dirMetaData{Path: templateRelPath, FileMode: info.Mode()})
		}
		return nil
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	e.TemplateEntries.Files = append(e.TemplateEntries.Files,
		fileMetaData{Path: templateRelPath, FileMode: info.Mode(), Content: espaceBackquote(string(content))})
	return nil
}

// As "Content" is quoted by "`", so we need to espace any "`" appears in "Content" via using `fmt.Sprintf`
func espaceBackquote(content string) string {
	if !strings.ContainsAny(content, "`") {
		return fmt.Sprintf("`%s`", content)
	}
	return fmt.Sprintf("fmt.Sprintf(`%s`, \"`\")", strings.ReplaceAll(content, "`", "%[1]s"))
}

func init() {
	outputTemplate = fmt.Sprintf(`// Code generated by go generate; DO NOT EDIT.
package {{ .InvokerPackage }}

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
	{{- range .TemplateEntries.Dirs }}
		dirMetaData{
			path: 		"{{ .Path }}",
			fileMode: 	{{ printf "%%d" .FileMode }},
		},
	{{- end }}
	}

	templateFiles = []fileMetaData{
	{{- range .TemplateEntries.Files }}
		fileMetaData{
			path: 		"{{ .Path }}",
			fileMode: 	{{ printf "%%d" .FileMode }},
			content:	%s,
		},
	{{- end }}
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
}`, "{{ .Content }}")
}
