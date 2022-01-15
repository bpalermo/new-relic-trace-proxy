
.PHONY: build
build:
	@docker build --progress plain --target export -t test . --output out
	@docker build --progress plain -t new-relic-trace-proxy .

.PHONY: test
test:
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...

