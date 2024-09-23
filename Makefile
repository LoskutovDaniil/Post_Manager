GO = $(shell which go)

.PHONY: all docker lint gen test

all:
	CGO_ENABLED=0 $(GO) build -mod=vendor -trimpath -ldflags="-w -s" -o ./dist/bin/ozon-test-task ./cmd/ozon-test-task/*.go

tooling/golangci-lint: vendor
	$(GO) build -mod=vendor -o ./tooling/golangci-lint ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint

tooling/gqlgen: vendor
	$(GO) build -mod=vendor -o ./tooling/gqlgen ./vendor/github.com/99designs/gqlgen

tooling/mockgen: vendor
	$(GO) build -mod=vendor -o ./tooling/mockgen ./vendor/go.uber.org/mock/mockgen

lint: tooling/golangci-lint
	$(PWD)/tooling/golangci-lint run --config ./golangci.yml

gen: tooling/gqlgen tooling/mockgen
	$(PWD)/tooling/gqlgen generate
	PATH="$(PWD)/tooling" $(GO) generate ./...

test:
	$(GO) test ./...
