# syntax=docker/dockerfile:1
FROM golang:1.17 AS build
WORKDIR /src
ARG TARGETOS
ARG TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -o /go/bin/app .

FROM build AS test
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
RUN mkdir -p /test
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go test -v -race -coverprofile=/test/coverage.out -covermode=atomic ./...

FROM scratch AS export
WORKDIR /
COPY --from=test /test/coverage.out .

FROM gcr.io/distroless/static@sha256:8ad6f3ec70dad966479b9fb48da991138c72ba969859098ec689d1450c2e6c97
COPY --from=build /go/bin/app /
CMD ["/app"]
