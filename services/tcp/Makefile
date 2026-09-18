SHELL := bash

COMMIT      ?= $(shell git rev-parse --short HEAD)
BUILD_TIME  ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
VERSION     := 1.x.x
SERVICE     := neosync.tcp

CMD_DIR     := ./cmd/server
BIN_DIR     := ./bin
CGO_ENABLED ?= 0


PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64

PLATFORMS_LIN := linux/amd64 linux/arm64

LDFLAGS := \
	-X 'main.version=$(VERSION)' \
	-X 'main.commit=$(COMMIT)' \
	-X 'main.buildTime=$(BUILD_TIME)'

.PHONY: run fmt lint test all build build-linux clean

all: build

build:
	@mkdir -p $(BIN_DIR)
	@echo ">> Building $(SERVICE) $(VERSION)"
	@for platform in $(PLATFORMS); do \
		OS=$${platform%/*}; \
		ARCH=$${platform#*/}; \
		EXT=""; \
		if [ "$$OS" = "windows" ]; then EXT=".exe"; fi; \
		OUT=$(BIN_DIR)/$(SERVICE)-$(VERSION)-$${OS}-$${ARCH}$${EXT}; \
		echo "  -> $$OUT"; \
		GOOS=$$OS \
		GOARCH=$$ARCH \
		CGO_ENABLED=$(CGO_ENABLED) \
		go build \
			-ldflags "$(LDFLAGS)" \
			-o $$OUT \
			$(CMD_DIR) || exit 1; \
	done

build-linux:
	@mkdir -p $(BIN_DIR)
	@echo ">> Building $(SERVICE) $(VERSION) (linux only)"
	@for platform in $(PLATFORMS_LIN); do \
		OS=$${platform%/*}; \
		ARCH=$${platform#*/}; \
		OUT=$(BIN_DIR)/$(SERVICE)-$(VERSION)-$${OS}-$${ARCH}; \
		echo "  -> $$OUT"; \
		GOOS=$$OS \
		GOARCH=$$ARCH \
		CGO_ENABLED=$(CGO_ENABLED) \
		go build \
			-ldflags "$(LDFLAGS)" \
			-o $$OUT \
			$(CMD_DIR) || exit 1; \
	done

HOST_OS   := $(shell go env GOOS)
HOST_ARCH := $(shell go env GOARCH)

ifeq ($(HOST_OS),windows)
	EXT := .exe
else
	EXT :=
endif

run:
	$(BIN_DIR)/$(SERVICE)-$(VERSION)-$(HOST_OS)-$(HOST_ARCH)$(EXT)

fmt:
	gofmt -w .
	goimports -w .

lint:
	golangci-lint run

check: fmt lint test
	@echo "==> All checks passed!"

test:
	@go test -v ./...

clean:
	rm -rf $(BIN_DIR)

