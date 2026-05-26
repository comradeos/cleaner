BINARY := build/cleaner
GOCACHE_DIR := $(CURDIR)/.gocache

.PHONY: help build run scan test fmt clean

help:
	@echo "Available targets:"
	@echo "  make build          Build the cleaner binary"
	@echo "  make run ARGS=...   Build and run the CLI"
	@echo "  make scan           Build and run 'cleaner scan'"
	@echo "  make test           Run Go tests"
	@echo "  make fmt            Format Go sources"
	@echo "  make clean          Remove build artifacts"

build:
	@mkdir -p build
	GOCACHE=$(GOCACHE_DIR) go build -o $(BINARY) .

run: build
	./$(BINARY) $(ARGS)

scan: build
	./$(BINARY) scan

test:
	GOCACHE=$(GOCACHE_DIR) go test ./...

fmt:
	gofmt -w main.go $$(find internal -name '*.go' -type f)

clean:
	rm -rf build bin .gocache cleaner
