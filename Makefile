.PHONY: run test

run:
	go run -race main.go

test:
	go test ./...
