.DEFAULT_GOAL := build

.PHONY: fmt vet staticcheck build
fmt:
	go fmt ./...

vet: fmt
	go vet ./...

staticcheck: vet
	staticcheck ./...

build: staticcheck
	go build 