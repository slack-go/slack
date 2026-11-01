.PHONY: help deps fmt vet lint staticcheck tidy-check test test-race test-integration pr-prep

help:
	@echo ""
	@echo "Welcome to slack-go/slack make."
	@echo "The following commands are available:"
	@echo ""
	@echo "    make deps              : Fetch all dependencies"
	@echo "    make fmt               : Format Go code and organize imports"
	@echo "    make vet               : Run go vet on all packages"
	@echo "    make lint              : Run golangci-lint for linting issues"
	@echo "    make staticcheck       : Run Staticcheck, including examples"
	@echo "    make tidy-check        : Check module files without modifying them"
	@echo "    make test              : Run all short tests"
	@echo "    make test-race         : Run all tests with race condition checking"
	@echo "    make test-integration  : Run all tests without limiting to short"
	@echo ""
	@echo "    make pr-prep           : Format code and run all PR checks"
	@echo ""

deps:
	@go mod tidy

fmt:
	@golangci-lint fmt

vet:
	@go vet ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0 run ./...

staticcheck:
	@staticcheck ./...

tidy-check:
	@go mod tidy -diff

test:
	@go test -v -count=1 -timeout 300s -short ./...

test-race:
	@go test -v -count=1 -timeout 300s -race ./...

test-integration:
	@go test -v -count=1 -timeout 600s ./...

pr-prep: fmt vet lint staticcheck test-race tidy-check
	@git diff --check
