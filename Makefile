.PHONY: tidy build core_build

tidy:
	go mod tidy

build:
	go build -o ./bin/anto-core.exe ./cmd/anto-core

core_build: build
