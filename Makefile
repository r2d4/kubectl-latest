GOOS ?= $(shell go env GOOS)
GOARCH = amd64
BUILD_DIR ?= ./out
ORG := github.com/r2d4
PROJECT := kubectl-latest
REPOPATH ?= $(ORG)/$(PROJECT)
BUILD_PKG := ./...

SUPPORTED_PLATFORMS := linux-$(GOARCH) darwin-$(GOARCH) windows-$(GOARCH).exe
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

$(BUILD_DIR)/$(PROJECT): $(BUILD_DIR)/$(PROJECT)-$(GOOS)-$(GOARCH)
	cp $(BUILD_DIR)/$(PROJECT)-$(GOOS)-$(GOARCH) $@

$(BUILD_DIR)/$(PROJECT)-%-$(GOARCH): $(GO_FILES) $(BUILD_DIR)
	GOOS=$* GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o $@ $(BUILD_PKG)

%.sha256: %
	shasum -a 256 $< > $@

%.exe: %
	cp $< $@

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

.PRECIOUS: $(foreach platform, $(SUPPORTED_PLATFORMS), $(BUILD_DIR)/$(PROJECT)-$(platform))
.PHONY: cross
cross: $(foreach platform, $(SUPPORTED_PLATFORMS), $(BUILD_DIR)/$(PROJECT)-$(platform).sha256)

.PHONY: install
install: $(GO_FILES) $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go install $(BUILD_PKG)

.SECONDEXPANSION:
TAR_TARGETS_linux   := $(BUILD_DIR)/$(PROJECT)-linux-amd64
TAR_TARGETS_darwin  := $(BUILD_DIR)/$(PROJECT)-darwin-amd64
TAR_TARGETS_windows := $(BUILD_DIR)/$(PROJECT)-windows-amd64.exe
$(BUILD_DIR)/$(PROJECT)-%-amd64.tar.gz: $$(TAR_TARGETS_$$*)
	tar -cvf $@ $^

.PHONY: cross-tars
cross-tars: $(BUILD_DIR)/$(PROJECT)-windows-amd64.tar.gz $(BUILD_DIR)/$(PROJECT)-linux-amd64.tar.gz $(BUILD_DIR)/$(PROJECT)-darwin-amd64.tar.gz
