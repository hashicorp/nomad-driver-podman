
.PHONY: changelogfmt
changelogfmt:
	@echo "--> Making [GH-xxxx] references clickable..."
	@sed -E 's|([^\[])\[GH-([0-9]+)\]|\1[[GH-\2](https://github.com/hashicorp/nomad-driver-podman/issues/\2)]|g' CHANGELOG.md > changelog.tmp && mv changelog.tmp CHANGELOG.md

.PHONY: check
check: ## Lint the source code
	@echo "==> Linting source code..."
	@golangci-lint run
	@echo "==> vetting hc-log statements"
	@hclogvet .

.PHONY: lint-deps
lint-deps: ## Install linter dependencies
## Keep versions in sync with tools/go.mod (see https://github.com/golang/go/issues/30515)
	@echo "==> Updating linter dependencies..."
	GO111MODULE=on cd tools && go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.24.0
	GO111MODULE=on cd tools && go get github.com/client9/misspell/cmd/misspell@v0.3.4
	GO111MODULE=on cd tools && go get github.com/hashicorp/go-hclog/hclogvet@master 
