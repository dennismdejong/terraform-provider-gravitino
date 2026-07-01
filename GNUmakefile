default: build

BINARY_NAME=terraform-provider-gravitino
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

build:
	go build -o bin/$(BINARY_NAME) -ldflags "-X main.version=$(VERSION)" .

test:
	go test -v -cover ./internal/...

testacc:
	TF_ACC=1 go test -v -cover ./internal/...

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/gravitino/gravitino/$(VERSION)/$$(go env GOOS)_$$(go env GOARCH)
	cp bin/$(BINARY_NAME) ~/.terraform.d/plugins/registry.terraform.io/gravitino/gravitino/$(VERSION)/$$(go env GOOS)_$$(go env GOARCH)/$(BINARY_NAME)

testacc-docker:
	docker compose run --rm test

generate:
	go generate ./...

.PHONY: build test testacc lint lint-fix fmt vet install testacc-docker generate
