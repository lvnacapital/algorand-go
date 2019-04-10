SRCPATH := $(GOPATH)/src/github.com/lvnacapital/algorand
TEST_SOURCES := $(shell cd $(SRCPATH) && go list ./...)

lint:
	golint `go list ./... | grep -v /vendor/`

build:
	cd $(SRCPATH) && go test -run xxx_phony_test $(TEST_SOURCES)
