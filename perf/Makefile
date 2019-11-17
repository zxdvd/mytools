
GO := go
GOBUILD := $(GO) build

all: kprobe execsnoop

linux: GOOS=linux
linux: GOARCH=amd64

linux: all

kprobe:
	$(GOBUILD) kprobe.go

execsnoop:
	$(GOBUILD) execsnoop.go

.PHONY: all linux kprobe execsnoop
