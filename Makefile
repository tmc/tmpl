.PHONY: all
all: install docs

.PHONY: install
install:
	go install

.PHONY: docs
docs:
	@cat README.in | sh -c "HELP='$$(tmpl -h 2>&1)' LATEST_TAG=$$(git tag -l |sort -V |tail -n1) LINUX_SHASUM=$$(grep linux_amd64 dist/checksums.txt |cut -d' ' -f1) MACOS_SHASUM=$$(grep darwin_amd64 dist/checksums.txt |cut -d' ' -f1) tmpl" > README.md

.PHONY: release
release:
	go run -mod=mod github.com/goreleaser/goreleaser@v1.9.2

.PHONY: release-dryrun
release-dryrun:
	go run -mod=mod github.com/goreleaser/goreleaser@v1.9.2 --snapshot --rm-dist
