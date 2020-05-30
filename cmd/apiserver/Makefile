.PHONY:build
build:
	go build -v ./cmd/apiserver

.PHONY: test
test:
	go test -v -rase -timeout 30s ./...
.DEFAULT_GOAL := build