# Spectra Makefile

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo "unknown")

LDFLAGS := -ldflags "\
	-X github.com/HarshalPatel1972/spectra/internal/version.Version=$(VERSION) \
	-X github.com/HarshalPatel1972/spectra/internal/version.Commit=$(COMMIT) \
	-X github.com/HarshalPatel1972/spectra/internal/version.Date=$(DATE)"

.PHONY: build test vet lint clean install

## build: Compile the spectra binary
build:
	go build $(LDFLAGS) -o bin/spectra ./cmd/spectra

## test: Run all tests with race detector
test:
	go test ./... -race -count=1

## vet: Run go vet
vet:
	go vet ./...

## lint: Run go vet (staticcheck can be added later)
lint: vet

## clean: Remove build artifacts
clean:
	rm -rf bin/ spectra-out/

## install: Install spectra to GOPATH/bin
install:
	go install $(LDFLAGS) ./cmd/spectra

## cover: Run tests with coverage
cover:
	go test ./... -race -count=1 -coverprofile=coverage.out
	go tool cover -func=coverage.out
