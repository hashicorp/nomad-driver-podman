
.PHONY: changelogfmt
changelogfmt:
	@echo "--> Making [GH-xxxx] references clickable..."
	@sed -E 's|([^\[])\[GH-([0-9]+)\]|\1[[GH-\2](https://github.com/hashicorp/nomad-driver-podman/issues/\2)]|g' CHANGELOG.md > changelog.tmp && mv changelog.tmp CHANGELOG.md

.PHONY: check
check: ## Lint the source code
	@echo "==> Linting source code..."
	@$(CURDIR)/build/bin/golangci-lint run
	@echo "==> vetting hc-log statements"
	@$(CURDIR)/build/bin/hclogvet $(CURDIR)

.PHONY: deps
deps: ## Install dependencies
## Keep versions in sync with tools/go.mod (see https://github.com/golang/go/issues/30515)
	@echo "==> Installing dependencies..."
	cd tools && GOBIN=$(CURDIR)/build/bin go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.24.0
	cd tools && GOBIN=$(CURDIR)/build/bin go get github.com/client9/misspell/cmd/misspell@v0.3.4
	cd tools && GOBIN=$(CURDIR)/build/bin go get github.com/hashicorp/go-hclog/hclogvet@master 
	cd tools && GOBIN=$(CURDIR)/build/bin go get gotest.tools/gotestsum
