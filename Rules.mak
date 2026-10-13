# MFP - comprehensive toolkit for multi-function printers and scanners
#
# Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
# See LICENSE for license terms and conditions
#
# Common part for Makefiles

# ----- Environment -----

TOPDIR  := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
GOFILES := $(wildcard *.go)
GOTESTS := $(wildcard *_test.go)
PACKAGE	:= $(shell basename $(shell pwd))

# ----- Parameters -----

GO	:= go
CTAGS	:= $(shell which ctags 2>/dev/null)
GOLINT	:= $(shell which golint 2>/dev/null)

# ----- Tools -----
IPP_REGISTRATIONS_GENERATE := go run $(TOPDIR)/tools/ipp-registrations-generate

# ----- Common targets -----

.PHONY: all
.PHONY: clean
.PHONY: cover
.PHONY: tags
.PHONY: test
.PHONY: vet

# Recursive targets
all:	do_all subdirs_all
clean:	do_clean subdirs_clean
lint:	do_lint subdirs_lint
test:	do_test subdirs_test
vet:	do_vet subdirs_vet

# Non-recursive targets
cover:	do_cover

# Dependencies
do_all:	tags $(ALL_LOCAL)

# Default actions
tags:
do_all:
do_clean:
do_cover:
do_lint:
do_test:
do_vet:

# Conditional actions
ifneq   ($(GOFILES),)

do_all:
	GOPROXY=off $(GO) build
ifneq	($(GOTESTS),)
	GOPROXY=off $(GO) test -c
	rm -f $(PACKAGE).test
endif

do_cover:
ifneq	($(GOTESTS),)
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm -f coverage.out
endif

do_lint:
ifneq	($(GOLINT),)
	$(GOLINT) -set_exit_status
endif

do_test:
ifneq	($(GOTESTS),)
	GOPROXY=off $(GO) test
endif

do_vet:
	GOPROXY=off $(GO) vet

endif

ifneq	($(CTAGS),)
tags:
	cd $(TOPDIR); rm -f tags; $(CTAGS) -R
endif

ifneq	($(CLEAN),)
do_clean:
	rm -f $(CLEAN)
endif

# ----- Subdirs handling

subdirs_all subdirs_lint subdirs_test subdirs_vet subdirs_clean:
	@for i in $(SUBDIRS); do \
                $(MAKE) -C $$i $(subst subdirs_,,$@) || exit 1; \
        done

