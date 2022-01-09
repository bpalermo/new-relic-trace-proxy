# syntax=docker/dockerfile:1
ARG GO_VERSION=1.17
FROM golang:${GO_VERSION} as base
WORKDIR /go/src/app
COPY go.* .
RUN go mod download
COPY . .

FROM base AS build
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o /go/bin/app .

FROM base AS test
RUN --mount=type=cache,target=/root/.cache/go-build \
    go test -v -race -coverprofile=coverage.out -covermode=atomic

FROM gcr.io/distroless/static@sha256:8ad6f3ec70dad966479b9fb48da991138c72ba969859098ec689d1450c2e6c97
COPY --from=build /go/bin/app /
CMD ["/app"]
