## Description

<!-- Brief description of changes -->

## Linked Issues

<!-- Resolves #123 -->

## Testing

- [ ] `go build ./...` passes
- [ ] `go test -v -cover ./internal/...` passes
- [ ] `golangci-lint run --config .github/golangci.yml ./...` passes
- [ ] `make generate` (if schema changed)

## Documentation

- [ ] `make generate` run (if schema changed)
- [ ] Examples updated in `examples/` directory

## Checklist

- [ ] Error handling uses `client.NewResourceError`
- [ ] 404 checks use `client.IsNotFoundError`
- [ ] `tflog.Debug` at CRUD boundaries
- [ ] Import with dot-separated ID parsing
- [ ] Enum fields have `stringvalidator.OneOf` validators
- [ ] Unit tests: schema, create, delete, import (valid + invalid)
- [ ] Resource/data source registered in `internal/provider/provider.go`
