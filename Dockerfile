# syntax=docker/dockerfile:1
FROM golang:1.17 AS build
WORKDIR /src
ARG TARGETOS
ARG TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -o /go/bin/app ./cmd/proxy && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -o /go/bin/healthchecker ./cmd/healthchecker

FROM gcr.io/distroless/static@sha256:8ad6f3ec70dad966479b9fb48da991138c72ba969859098ec689d1450c2e6c97
WORKDIR /
COPY --from=build /go/bin/app /
COPY --from=build /go/bin/healthchecker /
HEALTHCHECK \
    --interval=10s \
    --start-period=5s \
    --timeout=3s \
    CMD ["/healthchecker"]
EXPOSE 9001
USER 1000
ENTRYPOINT ["/app"]
