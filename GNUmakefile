SHELL = bash

# Handle multi-path environments
GOPATH := $(shell go env GOPATH | cut -d: -f1)

# Respect $GOBIN if set in environment or via $GOENV file.
BIN := $(shell go env GOBIN)
ifndef BIN
BIN := $(GOPATH)/bin
endif
default: help

HELP_FORMAT="    \033[36m%-25s\033[0m %s\n"
.PHONY: help
help: ## Display this usage information
	@echo "Valid targets:"
	@grep -E '^[^ ]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		sort | \
		awk 'BEGIN {FS = ":.*?## "}; \
			{printf $(HELP_FORMAT), $$1, $$2}'
	@echo ""

.PHONY: changelogfmt
changelogfmt: ## Format changelog GitHub links
	@echo "--> Making [GH-xxxx] references clickable..."
	@sed -E 's|([^\[])\[GH-([0-9]+)\]|\1[[GH-\2](https://github.com/hashicorp/nomad-driver-podman/issues/\2)]|g' CHANGELOG.md > changelog.tmp && mv changelog.tmp CHANGELOG.md

.PHONY: check
check: hclfmt ## Lint the source code
	@echo "==> Linting source code ..."
	@$(BIN)/golangci-lint run
	@echo "==> vetting hc-log statements"
	@$(BIN)/hclogvet $(CURDIR)

.PHONY: hclfmt
hclfmt: ## Format HCL files with hclfmt
	@echo "--> Formatting HCL"
	@find . -name '.git' -prune \
	        -o \( -name '*.nomad' -o -name '*.hcl' -o -name '*.tf' \) \
	      -print0 | xargs -0 hclfmt -w
	@if (git status -s | grep -q -e '\.hcl$$' -e '\.nomad$$' -e '\.tf$$'); then echo The following HCL files are out of sync; git status -s | grep -e '\.hcl$$' -e '\.nomad$$' -e '\.tf$$'; exit 1; fi

.PHONY: deps
deps: ## Install build dependencies
	@echo "==> Installing build dependencies ..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2
	go install github.com/hashicorp/go-hclog/hclogvet@v0.1.6
	go install gotest.tools/gotestsum@v1.10.0
	go install github.com/hashicorp/hcl/v2/cmd/hclfmt@d0c4fa8b0bbc2e4eeccd1ed2a32c2089ed8c5cf1

.PHONY: clean
clean: ## Cleanup previous build
	@echo "==> Cleanup previous build"
	rm -f ./build/nomad-driver-podman

.PHONY: dev
dev: clean build/nomad-driver-podman ## Build the nomad-driver-podman plugin

build/nomad-driver-podman:
	@echo "==> Building driver plugin ..."
	mkdir -p build
	CGO_ENABLED=0 \
	go build -o build/nomad-driver-podman .

.PHONY: test
test: ## Run unit tests
	@echo "==> Running unit tests ..."
	go test -v -race ./...

.PHONY: test-ci
test-ci: ## Run unit tests in CI
	@echo "==> Running unit tests in CI ..."
	@$(BIN)/gotestsum --format=testname --rerun-fails=0 --packages=". ./api" -- \
		-cover \
		-timeout=10m \
		-count=1 \
		. ./api

.PHONY: version
version:
ifneq (,$(wildcard version/version_ent.go))
	@$(CURDIR)/scripts/version.sh version/version.go version/version_ent.go
else
	@$(CURDIR)/scripts/version.sh version/version.go version/version.go
endif

# CRT release compilation
dist/%/nomad-driver-podman: GO_OUT ?= $@
dist/%/nomad-driver-podman:
	@echo "==> RELEASE BUILD of $@ ..."
	CGO_ENABLED=0 \
	GOOS=linux GOARCH=$(lastword $(subst _, ,$*)) \
	go build -trimpath -o $(GO_OUT)

# CRT release packaging (zip only)
.PRECIOUS: dist/%/nomad-driver-podman
dist/%.zip: dist/%/nomad-driver-podman
	@echo "==> RELEASE PACKAGING of $@ ..."
	@cp LICENSE $(dir $<)LICENSE.txt
	zip -j $@ $(dir $<)*
