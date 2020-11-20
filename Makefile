.PHONY: all
all: install docs

.PHONY: install
install:
	go install

.PHONY: docs
docs:
	@cat README.in | sh -c "HELP='$$(tmpl -h 2>&1)' LATEST_TAG=$$(git tag -l |tail -n1) LINUX_SHASUM=$$(grep linux_amd64 dist/checksums.txt |cut -d' ' -f1) tmpl" > README.md

.PHONY: release
release:
	goreleaser
