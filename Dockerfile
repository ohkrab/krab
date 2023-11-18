FROM golang:1.21.4-alpine3.18 AS build

LABEL org.opencontainers.image.source https://github.com/ohkrab/krab

WORKDIR /src
COPY go.* ./
RUN go mod download
RUN apk add --no-cache make

ENV CGO_ENABLED=0 \
    GOOS=linux    \
    GOARCH=amd64

ARG BUILD_VERSION=
ARG BUILD_DATE=
ARG BUILD_COMMIT=

COPY . ./
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN go build \
  -ldflags="-s -w -X 'github.com/ohkrab/krab/krab.InfoVersion=$BUILD_VERSION' -X 'github.com/ohkrab/krab/krab.InfoCommit=$BUILD_COMMIT' -X 'github.com/ohkrab/krab/krab.InfoBuildDate=$BUILD_DATE'" \
  -o /tmp/krab .

FROM alpine:3.18
COPY --from=build /tmp/krab /usr/local/bin/krab
ENTRYPOINT ["/usr/local/bin/krab"]

RUN mkdir -p /etc/krab

ENV KRAB_DIR=/etc/krab
