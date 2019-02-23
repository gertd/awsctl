# text reset
NO_COLOR=\033[0m
# green
OK_COLOR=\033[32;01m
#red
ERROR_COLOR=\033[31;01m
# cyan
WARN_COLOR=\033[36;01m
# yellow
ATTN_COLOR=\033[33;01m

PKGS := $(shell go list ./... | grep -v /vendor)

ROOT_DIR := $(git rev-parse --show-toplevel)
BIN_DIR  := $(PWD)/bin

GOMETALINTER := $(BIN_DIR)/gometalinter

GOOS :=
ifeq ($(OS),Windows_NT)
	GOOS = windows
else 
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		GOOS = linux
	endif
	ifeq ($(UNAME_S),Darwin)
		GOOS = darwin
	endif
endif
GOARCH ?= amd64

VERSION:=`git describe --tags --dirty`
COMMIT :=`git rev-parse --short HEAD 2>/dev/null`
DATE   :=`date "+%FT%T%z"`

LDFLAGS := -ldflags "-w -s -X github.com/gertd/awsctl/cmd.version=${VERSION} -X github.com/gertd/awsctl/cmd.date=${DATE} -X github.com/gertd/awsctl/cmd.commit=${COMMIT}"

BINARY := awsctl
VERSION ?= vlatest
PLATFORMS := windows linux darwin
OS = $(word 1, $@)

.PHONY: all
all: build test lint

.PHONY: build
build:
	@echo -e "$(ATTN_COLOR)==> build GOOS=$(GOOS) GOARCH=$(GOARCH) VERSION=$(VERSION)@$(COMMIT) $(NO_COLOR)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) GO111MODULE=on go build $(LDFLAGS) -o $(BIN_DIR)/aws-ctl ./

.PHONY: test
test:
	@echo -e "$(ATTN_COLOR)==> test $(NO_COLOR)"
	@go test $(PKGS)

$(GOMETALINTER):
	@echo -e "$(ATTN_COLOR)==> get gometalinter $(NO_COLOR)"
	@go get -u github.com/alecthomas/gometalinter
	@gometalinter --install 

.PHONY: lint
lint: $(GOMETALINTER)
	@echo -e "$(ATTN_COLOR)==> lint$(NO_COLOR)"
	@gometalinter ./... --vendor --deadline=90s
	@echo "$(NO_COLOR)\c"

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	@echo -e "$(ATTN_COLOR)==> release GOOS=$(GOOS) GOARCH=$(GOARCH) release/$(BINARY)-$(OS)-$(GOARCH) $(NO_COLOR)"
	@mkdir -p release
	@GOOS=$(OS) GOARCH=$(GOARCH) GO111MODULE=on go build $(LDFLAGS) -o release/$(BINARY)-$(OS)-$(GOARCH)$(if $(findstring $(OS),windows),".exe","")

.PHONY: release
release: windows linux darwin

.PHONY: install
install:
	@echo -e "$(ATTN_COLOR)==> install $(NO_COLOR)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) GO111MODULE=on go install $(LDFLAGS) ./
