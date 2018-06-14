help:
	@echo "build             run build"

.PHONY: build
build:
	gofmt -s -w .
	golint .
	go vet
	go build