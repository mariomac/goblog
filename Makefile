ASSETNAME   := $(shell basename $(shell pwd))
BINARY_NAME  = $(ASSETNAME)
GOCMD       ?= go
all: build

build: clean fmt lint test compile

clean:
	@echo "=== $(ASSETNAME) === [ clean ]: Removing binaries and coverage file..."
	@rm -rfv bin coverage.xml

lint:
	@echo "=== $(ASSETNAME) === [ lint ]: Validating source code running golint..."
	golangci-lint run

compile:
	@echo "=== $(ASSETNAME) === [ compile ]: Building $(BINARY_NAME)..."
	$(GOCMD) build -o bin/$(BINARY_NAME) ./src

test:
	@echo "=== $(ASSETNAME) === [ test ]: Running unit tests..."
	@gocov test ./src/... | gocov-xml > coverage.xml

fmt:
	@echo "=== $(ASSETNAME) === [ fmt ]: formatting code..."
	goimports -w ./src/

sample: compile
	@echo "=== $(ASSETNAME) === [ sample ]: running sample blog..."
	bin/$(BINARY_NAME) -cfg sample/config.yml


.PHONY: all build clean lint compile test fmt