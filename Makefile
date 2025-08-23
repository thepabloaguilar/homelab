.PHONY: test
test:
	@go test -v -race -vet=all -count=1 -coverprofile=coverage.out ./...

.PHONY: lint
lint:
	@go tool -modfile=.tools/go.mod golangci-lint run ./...

.PHONY: goimports
goimports:
	@go tool -modfile=.tools/go.mod goimports -local "$(shell go list -m)" -w .

.PHONY: goreleaser/check
goreleaser/check:
	@go tool -modfile=.tools/go.mod goreleaser check
