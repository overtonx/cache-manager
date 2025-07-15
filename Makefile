GO := go
GOLANGCI_LINT_VERSION := v1.59.1
GOLANGCI_LINT := $(shell which golangci-lint || echo "")
PROJECT_NAME := cache

.PHONY: lint test

all: lint test

lint:
	golangci-lint run --config .golangci.yml ./...

test:
	$(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

ci: fmt vet lint test