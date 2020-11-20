# tmpl

Command tmpl renders a template with the current env vars as input.

tmpl packs a punch in under 200 lines of code: a single static binary supplies the capabilities of
many more cmplicating templating engines.

It's especially helpful as an early entrypoint into containers to prepare configuration files.

```sh
$ tmpl -h
Usage of tmpl:
  -f string
    	Input source (default "-")
  -html
    	If true, use html/template instead of text/template
  -r string
    	If provided, traverse the argument as a directory, output is a tarball
  -w string
    	Output destination (default "-")
```

It includes all of the template helpers from [sprig](https://godoc.org/github.com/Masterminds/sprig).

It effectively exposes Go's [text/template](http://golang.org/pkg/text/template) for use in shells.

Reference [text/template](http://golang.org/pkg/text/template) documentation for template language specification.

## Safe Dockerfile Inclusion

To safely include in your build pipelines:
```Dockerfile
FROM ubuntu:bionic

RUN apt-get update
RUN apt-get install -y curl

ARG TMPL_URL=https://github.com/tmc/tmpl/releases/download/v1.8/tmpl_linux_amd64
ARG TMPL_SHA256SUM=9cb3a6d48405ef8bf7711d6e0be3a62e2c8257a147dcab1e04f6850f363eeed5
RUN curl -fsSLo tmpl ${TMPL_URL} \
		&& sha256sum tmpl \
		&& echo "${TMPL_SHA256SUM}  tmpl" | sha256sum -c - \
		&& chmod +x tmpl \
		&& mv tmpl /usr/local/bin/tmpl
```


### Example 1
Given a file 'a' with contents:


	{{ range $key, $value := . }}
	  KEY:{{ $key }} VALUE:{{ $value }}
	{{ end }}

Invoking

	$ cat a | env -i ANSWER=42 ITEM=Towel `which tmpl`

Produces


	KEY:ANSWER VALUE:42
	
	KEY:ITEM VALUE:Towel

### Example 2
Given a file 'b' with contents:


	VERSION={{.HEAD}}

Invoking


	$ cat b | HEAD="$(git rev-parse HEAD)" tmpl

Produces

	VERSION=4dce1b0a03b59b5d63c876143e9a9a0605855748

### Example 3
Given a directory via the `-r` flag, tmpl recurses, expanding each path and file and produces a tarball to the output destination.


Invoking

    $ mkdir testdata/recursive-out
	$ tmpl -r testdata/recursive-example | tar -C testdata/recursive-out --strip-components=2 -xvf -
	$ cat testdata/recursive-out/user-tmc

Produces (for me, at time of writing)

	For the current user tmc:
	Shell: /bin/bash
	EDITOR: vim
	😎
