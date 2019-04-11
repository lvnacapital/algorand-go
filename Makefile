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
# BUILDALL := $(SCRIPTDIR)/build.sh
# UPLOAD := $(SCRIPTDIR)/upload.sh
GO ?= go
GOX = gox
BUILD := $(GO) build
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
CAT := cat
AWK := awk
TR := tr
GIT := git
SYNC := aws s3 sync
SET := aws configure set
INVALIDATE := aws cloudfront create-invalidation

# Functions
# RWILDCARD = $(wildcard $1$2) $(foreach d,$(wildcard $1*),$(call RWILDCARD,$d/,$2))
HEAD := $(shell $(GIT) rev-parse --short HEAD | $(TR) -d "[ \r\n\']")
TAG := $(shell $(GIT) describe --always --tags --abbrev=0 | $(TR) -d "[v\r\n]")

# Files
# SHA256S = $(call RWILDCARD,$(BUILDDIR)/,*.sha256)
# BINARIES = $(patsubst $(BUILDDIR)/%.sha256,$(BUILDDIR)/%,$(SHA256S))
BINARIES := $(patsubst $(BUILDDIR)/$(WINDOWS)/%,$(BUILDDIR)/$(WINDOWS)/%.exe,$(addprefix $(BUILDDIR)/,$(addsuffix /$(BINARY),$(PLATFORMS))))
SHA256S := $(addsuffix .sha256,$(BINARIES))
VERIFY := $(patsubst $(BUILDDIR)/$(WINDOWS)/%,$(BUILDDIR)/$(WINDOWS)/%.exe,$(addsuffix /$(BINARY),$(PLATFORMS)))
GOPKG := Gopkg.lock

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

.PHONY: all clean deps lint test build buildall run verify upload sha256all verifyall $(VERIFY)

all: lint test run buildall sha256all verifyall

clean:
	@echo 'Cleaning...'
	$(CLEAN)
	$(RMDIR) $(BUILDDIR)

deps: $(GOPKG)
$(GOPKG):
	@echo 'Getting dependencies...'
	$(DEP) ensure

lint: deps
	@echo 'Linting...'
	$(LINT) -set_exit_status `$(LIST) ./... | $(GREP) -v '/vendor/'`

test: deps
	@echo 'Running tests...'
	$(TEST) -v ./...

build: $(BINARY)
$(BINARY): deps
	@echo 'Building single platform executable...'
#	$(MKDIR) $(BUILDDIR)
	$(BUILD) -o $(PWD)/$(BINARY) -v ./.
	$(CHMOD) +x $(PWD)/$(BINARY)

buildall: clean deps $(BINARIES)
$(BINARIES):
	@echo 'Building for all platforms...'
#	$(BUILDALL) -b $(BUILDDIR) -p $(PACKAGE) -o '$(PLATFORMS)'
	$(GOX) -ldflags="-s -X $(PACKAGE)/cmd.version=$(TAG) \
		-X $(PACKAGE)/cmd.commit=$(HEAD)" \
		-osarch "$(PLATFORMS)" -output="$(BUILDDIR)/{{.OS}}/{{.Arch}}/$(BINARY)"
	

sha256all: $(BINARIES) $(SHA256S)
$(SHA256S):
#	'$(@F)' is equivalent to '$(notdir $@)'
	@echo 'Generating SHA256 hash...'
	$(CAT) $(subst .sha256,,$@) | $(SHA256) | $(AWK) "{ print \$$1 \"  $(subst .sha256,,$(@F))\" }" > $@
#	@echo 'Verifying SHA256 checksum...'
#	printf $(dir $(patsubst $(BUILDDIR)/%.sha256,%,$@)) && cd $(dir $@) && $(SHA256) -c $(@F)

run: $(BINARY)
	@echo 'Checking single platform executable...'
	$(PWD)/$(BINARY)

verifyall: $(BINARIES) $(SHA256S) $(VERIFY)
$(VERIFY): %$(BINARY):
	@echo 'Verifying SHA256 checksums...'
	printf $(dir $*) && cd $(addprefix $(BUILDDIR)/,$(dir $*)) && $(SHA256) -c $(notdir $(wildcard $(BUILDDIR)/$(dir $*)*.sha256))

upload: $(BINARIES) $(SHA256S) verifyall
#	$(UPLOAD)
	@echo 'Uploading builds to AWS S3...'
	$(SYNC) $(BUILDDIR) $(BUCKET) --grants $(ALLACCESS) --region $(REGION)
#	@echo "Creating invalidation for AWS Cloudfront"
#	$(SET) preview.cloudfront true
#	$(INVALIDATE) --distribution-id $(DISTID) --paths /$(BINARY)