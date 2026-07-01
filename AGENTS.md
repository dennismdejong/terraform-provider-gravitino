# AGENTS.md

## Project: terraform-provider-gravitino

A Terraform provider for [Apache Gravitino](https://gravitino.apache.org/), built with the Plugin Framework (v1.19.0).

## Resources Implemented (all complete with CRUD + Import + tests)

| Resource              | TF Name              | Hierarchy                    |
|-----------------------|----------------------|------------------------------|
| Metalake              | `gravitino_metalake` | `metalake`                   |
| Catalog               | `gravitino_catalog`  | `metalake.catalog`           |
| Schema                | `gravitino_schema`   | `metalake.catalog.schema`    |
| Fileset               | `gravitino_fileset`  | `metalake.catalog.schema.fs` |
| Topic                 | `gravitino_topic`    | `metalake.catalog.schema.t`  |
| Table                 | `gravitino_table`    | `metalake.catalog.schema.tb` |
| View                  | `gravitino_view`     | (implemented)                |
| Model                 | `gravitino_model`    | (implemented)                |
| Function              | `gravitino_function` | (implemented)                |
| Partition             | `gravitino_partition`| (implemented)                |
| Tag                   | `gravitino_tag`      | (implemented)                |
| Policy                | `gravitino_policy`   | (implemented)                |
| Job                   | `gravitino_job`      | (implemented)                |

All resources also have corresponding data sources (list + get) and documentation under `docs/`.

## Architecture

```
main.go                        # Entry point
internal/
├── client/                    # HTTP client + per-resource REST methods
│   ├── client.go              # Core client: doRequest, Get, Post, Put, Delete
│   ├── metalake.go, catalog.go, schema.go, fileset.go, topic.go, ...
├── models/                     # JSON-annotated structs (requests & responses)
│   ├── common.go               # ErrorResponse, Audit, NameIdentifier
│   ├── metalake.go, catalog.go, schema.go, fileset.go, topic.go, ...
├── provider/
│   └── provider.go             # Provider schema + resource/data source registration
├── resources/                  # Terraform Plugin Framework resources
│   ├── metalake/resource.go    # Schema + CRUD + Import
│   ├── catalog/resource.go
│   ├── schema/resource.go
│   └── ...
├── datasources/                # Data source implementations
│   └── ...
```

## Conventions

### Patterns
- **ID format**: dot-separated hierarchy (`metalake`, `metalake.catalog`, etc.)
- **"id" attribute**: Always `Computed: true` with `stringplanmodifier.UseStateForUnknown()`
- **Audit**: `types.Object` with `AuditAttrTypes` (creator, create_time, last_modifier, last_modified_time)
- **Properties**: `types.Map` with `ElementType: types.StringType`
- **Update pattern**: Compare plan vs state for name/comment/properties, build `[]interface{}` update requests
- **404 handling**: Remove resource from state on 404/Not Found (except metalake)
- **Configure**: Casts `req.ProviderData` to `*client.Client`
- **Tests**: Use `httptest.NewServer` with custom handlers; test schema, create, delete, import

### Commands
- **Build**: `go build ./...`
- **Test (unit)**: `go test -v -cover ./internal/...`
- **Test (acceptance)**: `TF_ACC=1 go test -v -cover ./internal/...`
- **Docs**: Use `github.com/hashicorp/terraform-plugin-docs`

### Dependencies
- Go 1.26.4
- `terraform-plugin-framework` v1.19.0
- `terraform-plugin-framework-validators` v0.19.0
- `terraform-plugin-testing` v1.16.0
