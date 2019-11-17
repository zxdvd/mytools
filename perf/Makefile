
GO := go
GOBUILD := $(GO) build

all: kprobe execsnoop uprobe

linux:
	GOOS=linux GOARCH=amd64 make all

kprobe:
	$(GOBUILD) kprobe.go

execsnoop:
	$(GOBUILD) execsnoop.go

uprobe:
	$(GOBUILD) uprobe.go

.PHONY: all linux kprobe execsnoop uprobe
