.PHONY: build install test

build:
	mkdir -p bin/
	go build -o bin/krab main.go

install:
	cp bin/krab /usr/local/bin

test:
	go test -v ./...
