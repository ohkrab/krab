.PHONY: build install test

build:
	mkdir -p bin/
	go build -o bin/krab main.go

install:
	cp bin/krab /usr/local/bin

test:
	DATABASE_URL="postgres://krab:secret@localhost:5432/krab?sslmode=disable&prefer_simple_protocol=true" go test -v ./...

docker_test:
	docker run --rm -e DATABASE_URL="postgres://krab:secret@localhost:5432/krab?sslmode=disable" \
		-v ${HOME}/oh/krab/test/fixtures/simple:/etc/krab:ro qbart/krab:${BUILD_VERSION} version

docker_build:
	docker build -t qbart/krab:${BUILD_VERSION} \
		--build-arg BUILD_VERSION=${BUILD_VERSION} \
		--build-arg BUILD_COMMIT=${BUILD_COMMIT} \
		--build-arg BUILD_DATE=${BUILD_DATE} \
		.

docker_push:
	docker tag qbart/krab:${BUILD_VERSION} qbart/krab:latest
	docker push qbart/krab:${BUILD_VERSION}
	docker push qbart/krab:latest
