.PHONY: build install test

build:
	mkdir -p bin/
	go build -o bin/krab main.go

install:
	cp bin/krab /usr/local/bin

test:
	DATABASE_URL="postgres://krab:secret@localhost:5432/krab?sslmode=disable&prefer_simple_protocol=true" go test -v ./...
