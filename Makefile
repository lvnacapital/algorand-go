SHELL := /bin/bash
#.ONESHELL  # GNU make 3.82

# Configuration
USER := lvnacapital
BINARY := algorand
PACKAGE := github.com/$(USER)/$(BINARY)

# Directories
MAKEDIR := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
SRCDIR := $(PWD)
BUILDDIR := $(SRCDIR)/build
BINARY := algorand
BUILDDIR := $(PWD)/build
SCRIPTDIR := $(PWD)/scripts
LINUX := linux/amd64
WINDOWS := windows/amd64
DARWIN := darwin/amd64
PLATFORMS := $(LINUX) $(WINDOWS) $(DARWIN)

# Functions
RWILDCARD = $(wildcard $1$2) $(foreach d,$(wildcard $1*),$(call RWILDCARD,$d/,$2))
SHA256S = $(call RWILDCARD,$(BUILDDIR)/,*.sha256)
BINARIES = $(patsubst $(BUILDDIR)/%.sha256,$(BUILDDIR)/%,$(SHA256S))
VERIFY = printf $(dir $(patsubst $(BUILDDIR)/%.sha256,%,$(sha256))) && cd $(dir $(sha256)) && $(SHA256) -c $(sha256);

# Tools
GO ?= go
BUILD := $(GO) build
BUILDALL := $(SCRIPTDIR)/build.sh
CLEAN := $(GO) clean
TEST := $(GO) test
GET := $(GO) get -u
LINT := $(GO)lint
LIST := $(GO) list
DEP := dep
GREP := grep
WHICH := which
SHA256 := sha256sum
RMDIR := rm -rf
MKDIR := mkdir -p
CHMOD := chmod

# OS- and architecture-specific
ifeq ($(OS),Windows_NT)
    BINARY := $(patsubst %,%.exe,$(BINARY))
    ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
    else
        ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
        endif
        ifeq ($(PROCESSOR_ARCHITECTURE),x86)
        endif
    endif
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
    endif
    ifeq ($(UNAME_S),Darwin)
    endif
    UNAME_P := $(shell uname -p)
    ifeq ($(UNAME_P),x86_64)
    endif
    ifneq ($(filter %86,$(UNAME_P)),)
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
    endif
endif

.PHONY: all clean deps lint test build run verify

all: lint test build run verify

clean:
	@echo 'Cleaning...'
	$(CLEAN)
	$(RMDIR) $(BUILDDIR)
	
deps:
	@echo 'Getting dependencies...'
	$(DEP) ensure

lint: clean deps
	@echo 'Linting...'
	$(LINT) `$(LIST) ./... | $(GREP) -v /vendor/`

test: clean deps
	@echo 'Running tests...'
	$(TEST) -v ./...

$(BINARIES): build

build: clean deps
	@echo 'Building for all platforms...'
	$(BUILDALL) -b $(BUILDDIR) -p $(PACKAGE) -o '$(PLATFORMS)'

run: clean deps
	@echo 'Checking single platform executable...'
	# $(MKDIR) $(BUILDDIR)
	$(BUILD) -o $(PWD)/$(BINARY) -v ./.
	$(CHMOD) +x $(PWD)/$(BINARY)
	$(PWD)/$(BINARY)

verify: $(BINARIES)
	@echo 'Verifying SHA256 checksums...'
	@$(foreach sha256,$(SHA256S),$(VERIFY))
