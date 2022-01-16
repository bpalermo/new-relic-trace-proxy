
.PHONY: build
build:
	@docker build --progress plain --target export -t test . --output out
	@docker build --progress plain -t new-relic-trace-proxy .

.PHONY: test
test:
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...


.PHONY: lint
lint:
	@go vet ./...
	@golangci-lint run

.PHONY: install
install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0
