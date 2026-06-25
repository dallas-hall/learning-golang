.DEFAULT_GOAL := build
.PHONY:fmt vet build
fmt:
	go fmt ./...
vet: fmt
	go vet ./...
clean: vet
	go clean -i -x
build: clean
	go build