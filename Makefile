# Project variables
ORG         := khos2ow
NAME        := ratelimiter
DESCRIPTION := Example distributed Rate Limiter service
URL         := https://github.com/khos2ow/ratelimiter

# Repository variables
REPOSITORY  := github.com
NAMESPACE   := $(ORG)/$(NAME)
PACKAGE     := $(REPOSITORY)/$(NAMESPACE)

# Build variables
BUILD_DIR    := bin
COMMIT_HASH  ?= $(shell git rev-parse --short HEAD 2>/dev/null)
GIT_VERSION  ?= $(shell git describe --tags --exact-match 2>/dev/null || git describe --tags 2>/dev/null || echo "v0.0.1-$(COMMIT_HASH)")
BUILD_DATE   ?= $(shell date +%FT%T%z)
COVERAGE_OUT := coverage.out

# Go variables
GOOS        ?= $(shell go env GOOS)
GOARCH      ?= $(shell go env GOARCH)
GOPKGS      ?= $(shell go list $(MODVENDOR) ./... | grep -v /vendor)
GOFILES     ?= $(shell find . -type f -name '*.go' -not -path "./vendor/*")

GOLDFLAGS   :="
GOLDFLAGS   += -X $(PACKAGE)/internal/version.version=$(GIT_VERSION)
GOLDFLAGS   += -X $(PACKAGE)/internal/version.commitHash=$(COMMIT_HASH)
GOLDFLAGS   += -X $(PACKAGE)/internal/version.buildDate=$(BUILD_DATE)
GOLDFLAGS   +="

GOBUILD     ?= GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags $(GOLDFLAGS)
GORUN       ?= GOOS=$(GOOS) GOARCH=$(GOARCH) go run $(MODVENDOR)

# Docker variables
DEFAULT_TAG  ?= $(shell echo "$(GIT_VERSION)" | tr -d 'v')
DOCKER_IMAGE := $(NAMESPACE)
DOCKER_TAG   ?= $(DEFAULT_TAG)

# Binary versions
GITCHGLOG_VERSION := 0.9.1
GOLANGCI_VERSION  := v1.23.7

.PHONY: all
all: clean tools verify checkfmt lint test build

.PHONY: info
info: ## Show information about plugin
	@ echo "$(NAME) - $(GIT_VERSION) - $(BUILD_DATE)"

.PHONY: version
version: ## Show version of plugin
	@ echo "$(GIT_VERSION)"

#########################
## Development targets ##
#########################
.PHONY: clean
clean: ## Clean builds
	@ $(MAKE) --no-print-directory log-$@
	rm -rf ./$(BUILD_DIR) $(NAME) $(COVERAGE_OUT)

.PHONY: verify
verify: ## Verify 'vendor' dependencies
	@ $(MAKE) --no-print-directory log-$@
	go mod verify

.PHONY: tidy
tidy: ## Tidy up 'vendor' dependencies
	@ $(MAKE) --no-print-directory log-$@
	go mod tidy

.PHONY: lint
lint: ## Run linter
	@ $(MAKE) --no-print-directory log-$@
	golangci-lint run ./...

.PHONY: fmt
fmt: ## Format all go files
	@ $(MAKE) --no-print-directory log-$@
	goimports -w $(GOFILES)

.PHONY: checkfmt
checkfmt: RESULT ?= $(shell goimports -l $(GOFILES) | tee >(if [ "$$(wc -l)" = 0 ]; then echo "OK"; fi))
checkfmt: SHELL  := /usr/bin/env bash
checkfmt: ## Check formatting of all go files
	@ $(MAKE) --no-print-directory log-$@
	@ echo "$(RESULT)"
	@ if [ "$(RESULT)" != "OK" ]; then exit 1; fi

.PHONY: test
test: ## Run tests
	@ $(MAKE) --no-print-directory log-$@
	go test -coverprofile=$(COVERAGE_OUT) -covermode=atomic -v $(GOPKGS)

# removed and gitignoreed 'vendor/', not needed anymore #
.PHONY: vendor
vendor:

###################
## Build targets ##
###################
.PHONY: build
build: clean ## Build binary for current OS/ARCH
	@ $(MAKE) --no-print-directory log-$@
	$(GOBUILD) -o ./$(BUILD_DIR)/$(GOOS)-$(GOARCH)/$(NAME) $(PACKAGE)

.PHONY: build-all
build-all: GOOS   = linux darwin windows freebsd
build-all: GOARCH = amd64 arm
build-all: clean ## Build binary for all OS/ARCH
	@ $(MAKE) --no-print-directory log-$@
	@ ./scripts/build/build-all-osarch.sh "$(BUILD_DIR)" "$(NAME)" "$(GIT_VERSION)" "$(GOOS)" "$(GOARCH)" $(GOLDFLAGS)

.PHONY: docker
docker: ## Build Docker image
	@ $(MAKE) --no-print-directory log-$@
	docker build --pull --tag $(DOCKER_IMAGE):$(DOCKER_TAG) --file Dockerfile .

.PHONY: push
push: ## Push Docker image
	@ $(MAKE) --no-print-directory log-$@
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

####################
## Helper targets ##
####################
.PHONY: goimports
goimports:
ifeq (, $(shell which goimports))
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
endif

.PHONY: golangci
golangci:
ifeq (, $(shell which golangci-lint))
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s  -- -b $(shell go env GOPATH)/bin $(GOLANGCI_VERSION)
endif

.PHONY: gox
gox:
ifeq (, $(shell which gox))
	GO111MODULE=off go get -u github.com/mitchellh/gox
endif

.PHONY: tools
tools: ## Install required tools
	@ $(MAKE) --no-print-directory log-$@
	@ $(MAKE) --no-print-directory goimports golangci gox

#####################
## Release targets ##
#####################
.PHONY: release patch minor major
PATTERN =

release: VERSION ?= $(shell echo $(GIT_VERSION) | sed 's/^v//' | awk -F'[ .]' '{print $(PATTERN)}')
release: PUSH    := false
release: ## Prepare release
	@ $(MAKE) --no-print-directory log-$@
	@ ./scripts/release/release.sh "$(VERSION)" "$(PUSH)" "$(GIT_VERSION)" "1"

patch: PATTERN = '\$$1\".\"\$$2\".\"\$$3+1'
patch: release ## Prepare Patch release

minor: PATTERN = '\$$1\".\"\$$2+1\".0\"'
minor: release ## Prepare Minor release

major: PATTERN = '\$$1+1\".0.0\"'
major: release ## Prepare Major release

####################################
## Self-Documenting Makefile Help ##
####################################
.PHONY: help
help:
	@ grep -h -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

log-%:
	@ grep -h -E '^$*:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m==> %s\033[0m\n", $$2}'
