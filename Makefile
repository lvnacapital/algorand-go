# Configuration
USER := lvnacapital
BINARY := algorand
PACKAGE := github.com/$(USER)/$(BINARY)
BUCKET := s3://$(USER)/$(BINARY)
ALLACCESS := read=uri=http://acs.amazonaws.com/groups/global/AllUsers
REGION := us-west-2
DISTID := E3B5Z3LYG19QSL

# Directories
MAKEDIR := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
SRCDIR := $(PWD)
BUILDDIR := build
SCRIPTDIR := scripts
LINUX := linux/amd64
WINDOWS := windows/amd64
DARWIN := darwin/amd64
PLATFORMS := $(LINUX) $(WINDOWS) $(DARWIN)

# Tools
GO ?= go
GOX = gox
BUILD := $(GO) build
CLEAN := $(GO) clean
TEST := $(GO) test
GET := $(GO) get -u
FMT := $(GO)returns
LINT := $(GO)lint
LIST := $(GO) list
VET := $(GO) vet
DEP := dep
GREP := grep
WHICH := which
SHA256 := sha256sum
RMDIR := rm -rf
MKDIR := mkdir -p
CHMOD := chmod
CAT := cat
AWK := awk
TR := tr
GIT := git
SYNC := aws s3 sync
SET := aws configure set
INVALIDATE := aws cloudfront create-invalidation

# Functions
HEAD := $(shell $(GIT) rev-parse --short HEAD | $(TR) -d "[ \r\n\']")
TAG := $(shell $(GIT) describe --always --tags --abbrev=0 | $(TR) -d "[v\r\n]")
LDFLAGS := -s -X $(PACKAGE)/cmd.version=$(TAG) -X $(PACKAGE)/cmd.commit=$(HEAD)
OUTPUT := $(BUILDDIR)/{{.OS}}/{{.Arch}}/$(BINARY)

# Files
BINARIES := $(patsubst $(BUILDDIR)/$(WINDOWS)/%,$(BUILDDIR)/$(WINDOWS)/%.exe,$(addprefix $(BUILDDIR)/,$(addsuffix /$(BINARY),$(PLATFORMS))))
SHA256S := $(addsuffix .sha256,$(BINARIES))
VERIFY := $(patsubst $(BUILDDIR)/$(WINDOWS)/%,$(BUILDDIR)/$(WINDOWS)/%.exe,$(addsuffix /$(BINARY),$(PLATFORMS)))
PACKAGES = $(shell $(LIST) ./... | $(GREP) -v '/vendor/')
LOCK := Gopkg.lock

# OS- and architecture-specific
ifeq ($(OS),Windows_NT)
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

.PHONY: all clean deps fmt vet lint quicktest test build buildall run verify upload sha256all verifyall $(VERIFY)

all: fmt vet lint test run buildall sha256all verifyall

clean:
	@echo 'Cleaning...'
	$(CLEAN)
	$(RMDIR) $(BINARY) $(BUILDDIR)/

deps: $(LOCK)
$(LOCK):
	@echo 'Getting dependencies...'
	$(DEP) ensure

lint: deps
	@echo 'Linting...'
#	Capture output and force failure when there is non-empty output
	@echo '$(LINT) $(PACKAGES)'
	@OUTPUT=`$(LINT) $(PACKAGES) 2>&1`; \
	if [ "$$OUTPUT" ]; then \
		echo "'$(LINT)' errors:"; \
		echo "$$OUTPUT"; \
		exit 1; \
	fi

build: deps $(BINARY)
$(BINARY):
	@echo 'Building single platform executable...'
#	$(MKDIR) $(BUILDDIR)
	$(BUILD) -o $(BINARY) -ldflags='$(LDFLAGS)' -v ./.
	$(CHMOD) +x $(BINARY)

buildall: clean deps $(BINARIES)
$(BINARIES):
	@echo 'Building for all platforms...'
	$(GOX) -ldflags='$(LDFLAGS)' -osarch '$(PLATFORMS)' -output='$(OUTPUT)'

quicktest: deps
	@echo 'Running quick tests...'
	$(TEST) -v $(PACKAGES)

test: deps
	@echo 'Running tests...'
	$(TEST) -v -cover $(PACKAGES)

fmt: deps
	@echo 'Running format checks...'
	@echo "$(FMT) -l . | $(GREP) -v 'vendor[\/]'"
#	Capture output and force failure when there is non-empty output
	@OUTPUT=`$(FMT) -l . | $(GREP) -v 'vendor[\/]' 2>&1`; \
	if [ "$$OUTPUT" ]; then \
		echo "'$(FMT)' must be run on the following files:"; \
		echo "$$OUTPUT"; \
		exit 1; \
	fi

vet: deps
	go vet $(PACKAGES)

sha256all: $(BINARIES) $(SHA256S)
$(SHA256S):
#	'$(@F)' is equivalent to '$(notdir $@)'
	@echo 'Generating SHA256 hash...'
	$(CAT) $(subst .sha256,,$@) | $(SHA256) | $(AWK) "{ print \$$1 \"  $(subst .sha256,,$(@F))\" }" > $@
#	@echo 'Verifying SHA256 checksum...'
#	printf $(dir $(patsubst $(BUILDDIR)/%.sha256,%,$@)) && cd $(dir $@) && $(SHA256) -c $(@F)

run: $(BINARY)
	@echo 'Checking single platform executable...'
	./$(BINARY) --version

verifyall: $(BINARIES) $(SHA256S) $(VERIFY)
$(VERIFY): %$(BINARY):
	@echo 'Verifying SHA256 checksums...'
	printf $(dir $*) && cd $(addprefix $(BUILDDIR)/,$(dir $*)) && $(SHA256) -c $(notdir $(wildcard $(BUILDDIR)/$(dir $*)*.sha256))

upload: $(BINARIES) $(SHA256S) verifyall
	@echo 'Uploading builds to AWS S3...'
	$(SYNC) $(BUILDDIR) $(BUCKET) --grants $(ALLACCESS) --region $(REGION)
#	@echo "Creating invalidation for AWS Cloudfront"
#	$(SET) preview.cloudfront true
#	$(INVALIDATE) --distribution-id $(DISTID) --paths /$(BINARY)
