package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"text/template"

	"github.com/Masterminds/sprig"
)

var input = flag.String("f", "-", "Input source")

func main() {
	flag.Parse()
	in, err := getInput(*input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening input:", err)
		os.Exit(1)
	}
	if err := tmpl(in, os.Stdout, envMap()); err != nil {
		fmt.Fprintln(os.Stderr, "error opening input:", err)
		os.Exit(1)
	}
}

func getInput(path string) (io.Reader, error) {
	if path == "-" {
		return os.Stdin, nil
	}
	return os.Open(path)
}

func envMap() map[string]string {
	result := map[string]string{}
	for _, envvar := range os.Environ() {
		parts := strings.SplitN(envvar, "=", 2)
		result[parts[0]] = parts[1]
	}
	return result
}

func tmpl(in io.Reader, out io.Writer, ctx interface{}) error {
	i, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	tmpl, err := template.New("format string").Funcs(sprig.TxtFuncMap()).Parse(string(i))
	if err != nil {
		return err
	}
	return tmpl.Execute(out, ctx)
}
