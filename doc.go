// Command tmpl renders a template with the current env vars as input.
// # tmpl
// 
// Command tmpl renders a template with the current env vars as input.
// 
// tmpl packs a punch in under 200 lines of code: a single static binary supplies the capapbilities of
// many more cmplicating templating engines.
// 
// It's especially helpful as an early entrypoint into containers to prepare configuration files.
// 
// ```sh
// $ tmpl -h
// Usage of tmpl:
//   -f string
//     	Input source (default "-")
//   -html
//     	If true, use html/template instead of text/template
//   -r string
//     	If provided, traverse the argument as a directory, output is a tarball
//   -w string
//     	Output destination (default "-")
// ```
// 
// It includes all of the template helpers from [sprig](https://godoc.org/github.com/Masterminds/sprig).
// 
// It effectively exposes Go's [text/template](http://golang.org/pkg/text/template) for use in shells.
// 
// Reference [text/template](http://golang.org/pkg/text/template) documentation for template language specification.
// 
// ### Example 1
// Given a file 'a' with contents:
// 
// 
// 	{{ range $key, $value := . }}
// 	  KEY:{{ $key }} VALUE:{{ $value }}
// 	{{ end }}
// 
// Invoking
// 
// 	$ cat a | env -i ANSWER=42 ITEM=Towel `which tmpl`
// 
// Produces
// 
// 
// 	KEY:ANSWER VALUE:42
// 	
// 	KEY:ITEM VALUE:Towel
// 
// ### Example 2
// Given a file 'b' with contents:
// 
// 
// 	VERSION={{.HEAD}}
// 
// Invoking
// 
// 
// 	$ cat b | HEAD="$(git rev-parse HEAD)" tmpl
// 
// Produces
// 
// 
// 	VERSION=4dce1b0a03b59b5d63c876143e9a9a0605855748
package main
