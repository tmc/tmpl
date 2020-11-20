package main

import (
	"archive/tar"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	htmltemplate "html/template"
	"text/template"

	"github.com/Masterminds/sprig"
)

var (
	flagInput     = flag.String("f", "-", "Input source")
	flagOutput    = flag.String("w", "-", "Output destination")
	flagHTML      = flag.Bool("html", false, "If true, use html/template instead of text/template")
	flagRecursive = flag.String("r", "", "If provided, traverse the argument as a directory, output is a tarball")
)

func main() {
	flag.Parse()
	if err := run(*flagInput, *flagOutput, *flagRecursive, *flagHTML); err != nil {
		fmt.Fprintln(os.Stderr, "error opening input:", err)
		os.Exit(1)
	}
}

func run(input, output string, recurseDir string, htmlMode bool) error {
	in, err := getInput(*flagInput)
	if err != nil {
		return err
	}
	out, err := getOutput(*flagOutput)
	if err != nil {
		return err
	}
	if recurseDir != "" {
		return runDir(recurseDir, htmlMode, out, envMap())
	}
	return tmpl(in, *flagHTML, out, envMap())
}

func getInput(path string) (io.Reader, error) {
	if path == "-" {
		return os.Stdin, nil
	}
	return os.Open(path)
}

func getOutput(path string) (io.Writer, error) {
	if path == "-" {
		return os.Stdout, nil
	}
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
}

func envMap() map[string]string {
	result := map[string]string{}
	for _, envvar := range os.Environ() {
		parts := strings.SplitN(envvar, "=", 2)
		result[parts[0]] = parts[1]
	}
	return result
}

func tmpl(in io.Reader, htmlMode bool, out io.Writer, ctx interface{}) error {
	i, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	if htmlMode {
		tmpl, err := htmltemplate.New("format string").Funcs(sprig.HtmlFuncMap()).Parse(string(i))
		if err != nil {
			return err
		}
		return tmpl.Execute(out, ctx)
	}
	tmpl, err := template.New("format string").Funcs(sprig.TxtFuncMap()).Parse(string(i))
	if err != nil {
		return err
	}
	return tmpl.Execute(out, ctx)
}

func tmplToString(in io.Reader, htmlMode bool, ctx interface{}) (string, error) {
	o := bytes.NewBuffer([]byte{})
	err := tmpl(in, htmlMode, o, ctx)
	return o.String(), err
}

func tmplStr(in string, ctx interface{}) string {
	o := bytes.NewBuffer([]byte{})
	err := tmpl(strings.NewReader(in), false, o, ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "issue rendering string:", err)
	}
	return o.String()
}

func runDir(dir string, htmlMode bool, out io.Writer, ctx interface{}) error {
	tw := tar.NewWriter(out)
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		contents, err := tmplToString(f, htmlMode, ctx)
		if err != nil {
			return err
		}
		hdr := &tar.Header{
			Name: tmplStr(path, ctx),
			Mode: int64(info.Mode()),
			Size: int64(len(contents)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := tw.Write([]byte(contents)); err != nil {
			return err
		}
		return tw.Flush()
	})
}
