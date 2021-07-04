FROM golang:1.16-alpine AS build

WORKDIR /src
COPY go.* ./
RUN go mod download

ENV CGO_ENABLED=0 \
    GOOS=linux    \
    GOARCH=amd64

ARG BUILD_VERSION=
ARG BUILD_DATE=
ARG BUILD_COMMIT=

COPY . ./
RUN go build \
  -ldflags="-s -w -X 'github.com/ohkrab/krab/krab.InfoVersion=$BUILD_VERSION' -X 'github.com/ohkrab/krab/krab.InfoCommit=$BUILD_COMMIT' -X 'github.com/ohkrab/krab/krab.InfoBuildDate=$BUILD_DATE'" \
  -o /tmp/krab .

FROM alpine:3.14
COPY --from=build /tmp/krab /usr/local/bin/krab
ENTRYPOINT ["/usr/local/bin/krab"]

ENV KRAB_DIR=/etc/krab
ENV DATABASE_URL=

