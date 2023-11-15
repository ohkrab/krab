.PHONY: default build install test docker_test docker_build docker_push docker_nightly 

.PHONY: web
web:
	KRAB_AUTH=none \
	KRAB_AUTH_BASIC_USERNAME=krab \
	KRAB_AUTH_BASIC_PASSWORD=secret \
	DATABASE_URL="postgres://krab:secret@localhost:5432/krab?sslmode=disable" \
	air

.PHONY: gen
gen:
	templ generate

.PHONY: install
install:
	go install github.com/cosmtrek/air@latest
	go install github.com/a-h/templ/cmd/templ@latest

default:
	export DATABASE_URL="postgres://krab:secret@localhost:5432/krab?sslmode=disable" && \
	export KRAB_ENV=test && \
	export KRAB_DIR=./test/fixtures/tests && \
	make build && \
	./bin/krab test && \
	echo "ok"

build:
	mkdir -p bin/
	go build -o bin/krab main.go

test:
	DATABASE_URL="postgres://krab:secret@localhost:5432/krab?sslmode=disable&prefer_simple_protocol=true" go test -v ./... && echo "☑️ "

docker_test:
	docker run --rm -e DATABASE_URL="postgres://krab:secret@localhost:5432/krab?sslmode=disable" \
		-v ${HOME}/oh/krab/test/fixtures/simple:/etc/krab:ro ohkrab/krab-cli:${BUILD_VERSION} version

docker_build:
	docker build -t ohkrab/krab:${BUILD_VERSION} \
		--build-arg BUILD_VERSION=${BUILD_VERSION} \
		--build-arg BUILD_COMMIT=${BUILD_COMMIT} \
		--build-arg BUILD_DATE=${BUILD_DATE} \
		.

docker_push:
	docker tag ohkrab/krab:${BUILD_VERSION} ohkrab/krab:latest
	docker push ohkrab/krab:${BUILD_VERSION}
	docker push ohkrab/krab:latest

docker_nightly:
	docker build -t ohkrab/krab:nightly \
		--build-arg BUILD_VERSION=nightly \
		--build-arg BUILD_COMMIT=$$( git log -1 --pretty="format:%h" ) \
		--build-arg BUILD_DATE=$$( date -u +"%Y-%m-%dT%H:%M:%SZ" ) \
		.
	docker tag ohkrab/krab:nightly ohkrab/krab:latest
	docker push ohkrab/krab:nightly
