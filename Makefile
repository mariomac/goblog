ASSETNAME   := $(shell basename $(shell pwd))
BINARY_NAME  = $(ASSETNAME)
GOCMD       ?= go

GOLANGCI ?= go tool golangci-lint

all: build

build: clean fmt lint test compile

clean:
	@echo "=== $(ASSETNAME) === [ clean ]: Removing binaries and coverage file..."
	@rm -rfv bin coverage.xml

lint:
	@echo "=== $(ASSETNAME) === [ lint ]: Validating source code running golint..."
	$(GOLANGCI) run

compile:
	@echo "=== $(ASSETNAME) === [ compile ]: Building $(BINARY_NAME)..."
	$(GOCMD) build -o bin/$(BINARY_NAME) ./src

# TODO: coverage
test:
	@echo "=== $(ASSETNAME) === [ test ]: Running unit tests..."
	go test -race ./src/...

fmt:
	@echo "=== $(ASSETNAME) === [ fmt ]: formatting code..."
	$(GOLANGCI) fmt

sample: compile
	@echo "=== $(ASSETNAME) === [ sample ]: running sample blog..."
	bin/$(BINARY_NAME) -cfg sample/config.yml


.PHONY: all build clean lint compile test fmt