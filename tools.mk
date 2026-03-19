# Tools Makefile - Versioned Tool Management
# This file manages all external tools with specific versions for reproducible builds

# Platform detection
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)

ifeq ($(OS),darwin)
	PROTOC_OS := osx
else ifeq ($(OS),linux)
	PROTOC_OS := linux
else
	PROTOC_OS := $(OS)
endif

ifeq ($(ARCH),x86_64)
	PROTOC_ARCH := x86_64
else ifneq ($(filter arm64 aarch64,$(ARCH)),)
	PROTOC_ARCH := aarch_64
else
	PROTOC_ARCH := $(ARCH)
endif

# Tool versions
GOLANGCI_LINT_VERSION := v2.11.3
MOCKERY_VERSION := v3.7.0
PROTOC_VERSION := v34.0
PROTOC_GEN_GO_VERSION := v1.36.0
PROTOC_GEN_GO_GRPC_VERSION := v1.5.1

# Versioned tool paths
GOLANGCI_LINT := ./tools/golangci-lint-$(GOLANGCI_LINT_VERSION)
MOCKERY := ./tools/mockery-$(MOCKERY_VERSION)
PROTOC := ./tools/protoc-$(PROTOC_VERSION)
PROTOC_GEN_GO := ./tools/protoc-gen-go-$(PROTOC_GEN_GO_VERSION)
PROTOC_GEN_GO_GRPC := ./tools/protoc-gen-go-grpc-$(PROTOC_GEN_GO_GRPC_VERSION)

# Note: Using versioned binaries directly - no symlinks needed

.PHONY: tools
tools: $(GOLANGCI_LINT) $(MOCKERY) $(PROTOC) $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC)

.PHONY: clean-tools
clean-tools:
	rm -f ./tools/golangci-lint-*
	rm -f ./tools/mockery-*
	rm -f ./tools/protoc-*
	rm -f ./tools/protoc-gen-go-*
	rm -f ./tools/protoc-gen-go-grpc-*
	rm -rf ./tools/include


# golangci-lint installation with version
$(GOLANGCI_LINT):
	@mkdir -p ./tools
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./tools $(GOLANGCI_LINT_VERSION)
	mv ./tools/golangci-lint $@

# mockery installation with version
$(MOCKERY):
	@mkdir -p ./tools
	curl -sSfL https://github.com/vektra/mockery/releases/download/$(MOCKERY_VERSION)/mockery_$(subst v,,$(MOCKERY_VERSION))_$$(uname -s)_$$(uname -m).tar.gz | tar -xz -C ./tools mockery
	mv ./tools/mockery $@


# protoc installation with version
$(PROTOC):
	@mkdir -p ./tools
	curl -sSfL https://github.com/protocolbuffers/protobuf/releases/download/$(PROTOC_VERSION)/protoc-$(subst v,,$(PROTOC_VERSION))-$(PROTOC_OS)-$(PROTOC_ARCH).zip -o protoc.zip
	unzip protoc.zip
	mv bin/protoc $@
	mv include ./tools/include
	rm -rf bin include readme.txt protoc.zip

# protoc-gen-go installation with version
$(PROTOC_GEN_GO):
	@mkdir -p ./tools
	GOBIN=$(shell pwd)/tools go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	mv ./tools/protoc-gen-go $@
	ln -sf $(notdir $@) ./tools/protoc-gen-go

# protoc-gen-go-grpc installation with version
$(PROTOC_GEN_GO_GRPC):
	@mkdir -p ./tools
	GOBIN=$(shell pwd)/tools go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	mv ./tools/protoc-gen-go-grpc $@
	ln -sf $(notdir $@) ./tools/protoc-gen-go-grpc

# Help target
.PHONY: help-tools
help-tools:
	@echo "Available tool targets:"
	@echo "  tools        - Install all tools with specific versions"
	@echo "  clean-tools  - Remove all versioned tools"
	@echo ""
	@echo "Tool versions:"
	@echo "  golangci-lint:      $(GOLANGCI_LINT_VERSION)"
	@echo "  mockery:            $(MOCKERY_VERSION)"
	@echo "  protoc:             $(PROTOC_VERSION)"
	@echo "  protoc-gen-go:      $(PROTOC_GEN_GO_VERSION)"
	@echo "  protoc-gen-go-grpc: $(PROTOC_GEN_GO_GRPC_VERSION)"
