package main

import (
	"archive/tar"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	htmltemplate "html/template"
	"text/template"

	"github.com/tmc/tmpl/sprig"
)

var (
	flagInput     = flag.String("f", "-", "Input source")
	flagOutput    = flag.String("w", "-", "Output destination")
	flagHTML      = flag.Bool("html", false, "If true, use html/template instead of text/template")
	flagRecursive = flag.String("r", "", "If provided, traverse the argument as a directory")
	flagStripN    = flag.Int("stripn", 0, "If provided, strips this many directories from the output (only valid if -r and -w are provided)")
)

func main() {
	flag.Parse()
	if err := run(*flagInput, *flagOutput, *flagRecursive, *flagHTML); err != nil {
		fmt.Fprintln(os.Stderr, "tmpl error:", err)
		os.Exit(1)
	}
}

func run(input, output string, recurseDir string, htmlMode bool) error {
	in, err := getInput(input)
	if err != nil {
		return err
	}
	if recurseDir != "" {
		return runDir(recurseDir, htmlMode, output, *flagStripN, envMap())
	}
	out, err := getOutput(*flagOutput)
	if err != nil {
		return err
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
	i, err := io.ReadAll(in)
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

func runDir(dir string, htmlMode bool, outPath string, stripN int, ctx interface{}) error {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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
			return fmt.Errorf("%v: %w", path, err)
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
	if err != nil {
		return fmt.Errorf("issue recursing: %w", err)
	}
	if outPath == "-" {
		_, err = io.Copy(os.Stdout, buf)
		return err
	}
	return extractTar(buf, outPath, stripN)
}

func extractTar(buf io.Reader, outPath string, stripN int) error {
	tarReader := tar.NewReader(buf)
	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Next() failed: %w", err)
		}
		path := header.Name
		parts := strings.Split(path, string(filepath.Separator))
		toStrip := stripN
		if toStrip >= len(parts) {
			toStrip = len(parts) - 1
		}
		if len(parts) > toStrip {
			path = strings.Join(parts[toStrip:], string(filepath.Separator))
		}
		fullPath := filepath.Join(outPath, path)
		if err := ensureEnclosingDir(fullPath); err != nil {
			return fmt.Errorf("issue ensuring directory exists: %w", err)
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(fullPath, 0755); err != nil {
				return fmt.Errorf("Mkdir() failed: %w", err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(fullPath)
			if err != nil {
				return fmt.Errorf("Create() failed: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("io.Copy() failed: %w", err)
			}
			outFile.Close()

		default:
			return fmt.Errorf("extractTar: uknown type: %v in %v", header.Typeflag, header.Name)
		}
	}
	return nil
}

func ensureEnclosingDir(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0755)
}
