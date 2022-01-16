IMAGE_NAME:=new-relic-trace-proxy
IMAGE_TAG:=latest

.PHONY: build
build:
	@docker build --progress plain -t $(IMAGE_NAME):$(IMAGE_TAG) .

.PHONY: test
test:
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...

.PHONY: run
run:
	@docker run --rm -p 9001:9001 $(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: lint
lint:
	@go vet ./...
	@golangci-lint run

.PHONY: install
install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0
