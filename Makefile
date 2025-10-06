.PHONY: all
all: install docs

.PHONY: install
install:
	go install

.PHONY: docs
docs:
	@LATEST_TAG=$$(gh release list --exclude-pre-releases -L1 --jq '.[0].tagName'); \
	SHAS=$$(gh api repos/tmc/tmpl/releases/tags/$$LATEST_TAG --jq '.assets[] | select(.name | startswith("tmpl-") and contains("amd64") and (endswith(".jsonl")|not)) | .name + " " + .digest' | sed 's/sha256://'); \
	HELP="$$(tmpl -h 2>&1)" \
	LATEST_TAG="$$LATEST_TAG" \
	LINUX_SHASUM=$$(echo "$$SHAS" | awk '/linux/{print $$2}') \
	MACOS_SHASUM=$$(echo "$$SHAS" | awk '/darwin/{print $$2}') \
	tmpl < README.in > README.md
