SHELL := bash
GO := go

GO_SOURCES := $(shell find . -type f -name '*.go')

git-multiswitch: $(GO_SOURCES)
	$(GO) build

clean:
	rm -f git-multiswitch
