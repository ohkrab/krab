.PHONY: build install test

build:
	mkdir -p bin/
	go build -o bin/krab main.go

install:
	cp bin/krab /usr/local/bin

test:
	DATABASE_URL="postgres://krab:secret@localhost:5432/krab?sslmode=disable" go test -v ./...
