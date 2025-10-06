# tmpl

Command tmpl renders a template with the current env vars as input.

tmpl packs a punch in under 200 lines of code: a single static binary supplies the capabilities of
many more complicated templating engines.

It's especially helpful as an early entrypoint into containers to prepare configuration files.

```sh
$ tmpl -h
Usage of tmpl:
  -f string
    	Input source (default "-")
  -html
    	If true, use html/template instead of text/template
  -r string
    	If provided, traverse the argument as a directory
  -stripn int
    	If provided, strips this many directories from the output (only valid if -r and -w are provided)
  -w string
    	Output destination (default "-")
```

It effectively exposes Go's [text/template](http://golang.org/pkg/text/template) for use in shells.

Reference [text/template](http://golang.org/pkg/text/template) documentation for template language specification.

It includes all of the template helpers from [sprig](https://godoc.org/github.com/Masterminds/sprig).

## Safe Dockerfile Inclusion

To safely include in your build pipelines:
```Dockerfile
FROM ubuntu:bionic

RUN apt-get update
RUN apt-get install -y curl

ARG TMPL_URL=https://github.com/tmc/tmpl/releases/download/v1.12.4/tmpl_linux_amd64
ARG TMPL_SHA256SUM=13f00727845c5d9611b710165076328311ee8681991dcc2a8f4f10c81d73cbd7
RUN curl -fsSLo tmpl ${TMPL_URL} \
		&& echo "${TMPL_SHA256SUM}  tmpl" | sha256sum -c - \
		&& chmod +x tmpl && mv tmpl /usr/local/bin/tmpl
```

## Safe Shell Scripting Inclusion

To safely include in your shell scripts:
```bash
#!/bin/bash
set -euo pipefail

# Helper Functions
case "${OSTYPE}" in
linux*) platform=linux
	;;
darwin*)
	platform=darwin
	;;
*) platform=unknown ;;
esac

function install_tmpl() {
  if [[ "${platform}" == "darwin" ]]; then
    TMPL_SHA256SUM=c8931b6f36c21e443b87ea50d7fc28877025907c2376fbe5144327bcaf2c689c
  else
    TMPL_SHA256SUM=13f00727845c5d9611b710165076328311ee8681991dcc2a8f4f10c81d73cbd7
  fi
  TMPL_URL=https://github.com/tmc/tmpl/releases/download/v1.12.4/tmpl_${platform}_amd64
  curl -fsSLo tmpl ${TMPL_URL} \
    && echo "${TMPL_SHA256SUM}  tmpl" | sha256sum -c - \
    && chmod +x tmpl
  mv tmpl /usr/local/bin/tmpl || echo "could not move tmpl into place"
}

command -v tmpl > /dev/null || install_tmpl
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
