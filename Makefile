ASSETNAME  := $(shell basename $(shell pwd))
BINARY_NAME   = $(ASSETNAME)
GO_PKGS      := $(shell go list ./... | grep -v "/vendor/")

all: build

build: clean deps validate test compile

clean:
	@echo "=== $(ASSETNAME) === [ clean ]: Removing binaries and coverage file..."
	@rm -rfv bin coverage.xml

deps:
	@echo "=== $(ASSETNAME) === [ deps ]: Updating package dependencies required by the project..."
	glide update

validate:
	@echo "=== $(ASSETNAME) === [ validate ]: Validating source code running golint..."
	golint src/...

compile:
	@echo "=== $(ASSETNAME) === [ compile ]: Building $(BINARY_NAME)..."
	go build -o bin/$(BINARY_NAME) ./src

test:
	@echo "=== $(ASSETNAME) === [ test ]: Running unit tests..."
	@gocov test $(GO_PKGS) | gocov-xml > coverage.xml

.PHONY: all build clean deps validate compile test