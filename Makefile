.PHONY: all
all: install docs

.PHONY: install
install:
	go install

.PHONY: docs
docs:
	cat README.in | HELP=$(tmpl -h 2>&1) tmpl > README.md

.PHONY: release
release:
	goreleaser
