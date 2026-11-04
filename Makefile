.PHONY: help deps fmt lint tidy-check test test-race test-integration pr-prep

help:
	@echo ""
	@echo "Welcome to slack-go/slack make."
	@echo "The following commands are available:"
	@echo ""
	@echo "    make deps              : Fetch all dependencies"
	@echo "    make fmt               : Format Go code and organize imports"
	@echo "    make lint              : Run golangci-lint"
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

lint:
	@golangci-lint run ./...

tidy-check:
	@go mod tidy -diff

test:
	@go test -v -count=1 -timeout 300s -short ./...

test-race:
	@go test -v -count=1 -timeout 300s -race ./...

test-integration:
	@go test -v -count=1 -timeout 600s ./...

pr-prep: fmt lint test-race tidy-check
	@git diff --check
