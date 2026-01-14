# ===============================
# Project metadata
# ===============================
APP_NAME := zyra
MODULE   := github.com/Mahmoud-Khaled-FS/zyra

# ===============================
# Version info
# ===============================
VERSION ?= 0.0.1
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# ===============================
# Go settings
# ===============================
GO      := go
GOFLAGS := -trimpath

LDFLAGS := \
	-X $(MODULE)/internal/version.Version=$(VERSION) \
	-X $(MODULE)/internal/version.Commit=$(COMMIT) \
	-X $(MODULE)/internal/version.Date=$(DATE)

# ===============================
# Targets
# ===============================

.PHONY: all build run test clean version

all: build

build:
	@echo "Building $(APP_NAME)"
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(APP_NAME)

run:
	$(GO) run -ldflags "$(LDFLAGS)" ./main.go

version:
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Date:    $(DATE)"

clean:
	rm -f $(APP_NAME)
