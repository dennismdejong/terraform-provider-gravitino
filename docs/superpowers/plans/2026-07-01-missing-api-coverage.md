# Missing Gravitino API Features — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Close API coverage gaps — implement model versions, job template CRUD, statistics CRUD, lineage, catalog test, fileset files, and missing client methods.

**Architecture:** Each feature follows the standard 3-layer pattern: `models/*.go` → `client/*.go` → `resources/*/` + `datasources/*/`. Resources use the same conventions as existing code (`client.NewResourceError`, `client.IsNotFoundError`, `tflog.Debug`, `stringvalidator.OneOf`).

**Tech Stack:** Go 1.26.4, terraform-plugin-framework v1.19.0, terraform-plugin-log v0.10.0

**Current coverage:** 17 resources, 41 data sources, 130 REST endpoints — 24 endpoints missing

## Global Constraints

- `go build ./...` must pass after each task
- `go test -v -cover ./internal/...` must pass after each task
- All error handling via `client.NewResourceError`
- All 404 checks via `client.IsNotFoundError`
- All CRUD boundaries have `tflog.Debug`
- All enum fields have `stringvalidator.OneOf` with values from API docs
- Resources registered in `internal/provider/provider.go`
- Docs generated via `make generate`

---

## Amount of work

| Phase | Features | Endpoints | New files | Effort |
|-------|----------|-----------|-----------|--------|
| 1 | Model versions | 10 | ~8 | Large |
| 2 | Job templates | 4 | ~8 | Medium |
| 3 | Statistics + overige | 10 | ~12 | Medium |

---

### Task 1: Model versions — Models + Client

**Files:**
- Create: `internal/models/model_version.go`
- Create: `internal/client/model.go` (expand existing, or new `internal/client/model_version.go`)

**Interfaces:**
- Produces: Model version structs (ModelVersion, ModelVersionListResponse, ModelVersionCreateRequest, etc.)
- Produces: Client methods (ListModelVersions, LinkModelVersion, GetModelVersion, UpdateModelVersion, DeleteModelVersion, GetModelVersionByAlias, DeleteModelVersionByAlias, UpdateModelVersionByAlias, GetModelVersionURI, GetModelVersionURIByAlias)

**API Reference:** [Gravitino Model REST API](https://gravitino.apache.org/docs/next/api/rest/list-models)

#### Model structs needed:

```go
type ModelVersionCreateRequest struct {
    Version string            `json:"version"`
    URI     string            `json:"uri,omitempty"`
    Aliases []string          `json:"aliases,omitempty"`
    Comment string            `json:"comment,omitempty"`
    Properties map[string]string `json:"properties,omitempty"`
}

type ModelVersionResponse struct {
    Code          int32        `json:"code"`
    ModelVersion  ModelVersion `json:"modelVersion"`
}

type ModelVersion struct {
    Version    string            `json:"version"`
    URI        string            `json:"uri,omitempty"`
    Aliases    []string          `json:"aliases,omitempty"`
    Comment    string            `json:"comment,omitempty"`
    Properties map[string]string `json:"properties,omitempty"`
    Audit      *Audit            `json:"audit,omitempty"`
}

type ModelVersionListResponse struct {
    Code           int32          `json:"code"`
    ModelVersions  []ModelVersion `json:"modelVersions"`
}

type ModelVersionURIResponse struct {
    Code int32  `json:"code"`
    URI  string `json:"uri"`
}
```

#### Client methods:

| Method | HTTP | Path |
|--------|------|------|
| ListModelVersions(m, c, s, model) | GET | `/metalakes/:m/catalogs/:c/schemas/:s/models/:model/versions` |
| LinkModelVersion(m, c, s, model, req) | POST | `/metalakes/:m/catalogs/:c/schemas/:s/models/:model/versions` |
| GetModelVersion(m, c, s, model, version) | GET | `.../models/:model/versions/:version` |
| UpdateModelVersion(m, c, s, model, version, req) | PUT | Same |
| DeleteModelVersion(m, c, s, model, version) | DELETE | Same |
| GetModelVersionByAlias(m, c, s, model, alias) | GET | `.../models/:model/versions/aliases/:alias` |
| DeleteModelVersionByAlias(m, c, s, model, alias) | DELETE | Same |
| UpdateModelVersionByAlias(m, c, s, model, alias, req) | PUT | Same |
| GetModelVersionURI(m, c, s, model, version) | GET | `.../models/:model/versions/:version/uri` |
| GetModelVersionURIByAlias(m, c, s, model, alias) | GET | `.../models/:model/versions/aliases/:alias/uri` |

---

### Task 2: Model versions — Resource + Datasources

**Files:**
- Create: `internal/resources/model_version/resource.go`
- Create: `internal/resources/model_version/resource_test.go`
- Create: `internal/datasources/model_version/datasource_get.go`
- Create: `internal/datasources/model_version/datasource_list.go`
- Create: `internal/datasources/model_version/datasource_test.go`

**Terraform resource:** `gravitino_model_version`

Schema attributes:
- id: Computed, UseStateForUnknown (format: metalake.catalog.schema.model.version)
- metalake: Required
- catalog: Required
- schema: Required
- model: Required
- version: Required
- uri: Optional, Computed
- aliases: Optional, Computed (List of String)
- comment: Optional, Computed
- properties: Optional, Computed (Map of String)
- audit: Computed (Object)

Create: POST linkModelVersion
Read: GET GetModelVersion with 404 check
Update: Compare plan vs state for uri/aliases/comment/properties
Delete: DELETE DeleteModelVersion
Import: dot-separated "metalake.catalog.schema.model.version"

**Data sources:**
- `gravitino_model_version` — get by version or alias
- `gravitino_model_versions` — list all versions for a model

**Registration:** `rsmodelversion.NewModelVersionResource`, `dsmodelversion.NewModelVersionDataSource`, `dsmodelversion.NewModelVersionsDataSource`

---

### Task 3: Job templates — Models + Client

**Files:**
- Create/Expand: `internal/models/job.go` (add JobTemplate structs)
- Expand: `internal/client/job.go` (add 4 template methods)

**API Reference:** [Gravitino Job REST API](https://gravitino.apache.org/docs/next/api/rest/list-job-templates)

#### New structs:

```go
type JobTemplateCreateRequest struct {
    Name        string            `json:"name"`
    Template    string            `json:"template"`
    Parameters  map[string]string `json:"parameters,omitempty"`
    Comment     string            `json:"comment,omitempty"`
    Properties  map[string]string `json:"properties,omitempty"`
}

type JobTemplateResponse struct {
    Code         int32       `json:"code"`
    JobTemplate  JobTemplate `json:"jobTemplate"`
}

type JobTemplate struct {
    Name        string            `json:"name"`
    Template    string            `json:"template"`
    Parameters  map[string]string `json:"parameters,omitempty"`
    Comment     string            `json:"comment,omitempty"`
    Properties  map[string]string `json:"properties,omitempty"`
    Audit       *Audit            `json:"audit,omitempty"`
}

type JobTemplateListResponse struct {
    Code          int32         `json:"code"`
    JobTemplates  []JobTemplate `json:"jobTemplates"`
}
```

#### Client methods:

| Method | HTTP | Path |
|--------|------|------|
| RegisterJobTemplate(m, req) | POST | `/metalakes/:m/jobs/templates` |
| GetJobTemplate(m, name) | GET | `/metalakes/:m/jobs/templates/:name` |
| UpdateJobTemplate(m, name, req) | PUT | Same |
| DeleteJobTemplate(m, name) | DELETE | Same |

Note: `ListJobTemplates` already exists.

---

### Task 4: Job templates — Resource + Datasources

**Files:**
- Create: `internal/resources/job_template/resource.go`
- Create: `internal/resources/job_template/resource_test.go`
- Create: `internal/datasources/job_template/datasource_get.go`
- Create: `internal/datasources/job_template/datasource_list.go`
- Create: `internal/datasources/job_template/datasource_test.go`

**Terraform resource:** `gravitino_job_template`

Schema:
- id: Computed (format: metalake.template_name)
- metalake: Required
- name: Required
- template: Required
- parameters: Optional, Computed (Map of String)
- comment: Optional, Computed
- properties: Optional, Computed (Map of String)
- audit: Computed

**Data sources:** `gravitino_job_template`, `gravitino_job_templates`

---

### Task 5: Statistics CRUD — Resource

**Files:**
- Expand: `internal/client/statistics.go`
- Create: `internal/resources/statistics/resource.go`
- Create: `internal/resources/statistics/resource_test.go`

**API Reference:** [Gravitino Statistics API](https://gravitino.apache.org/docs/next/api/rest/list-statistics)

#### Client additions:

```go
func (c *Client) UpdateStatistics(metalake, objType, objName string, body interface{}) error
func (c *Client) DeleteStatistics(metalake, objType, objName string) error
```

#### Terraform resource additions:

Expand existing `gravitino_statistics` from read-only datasource to include a resource for updating/deleting statistics. The resource Create = PUT update-statistics, Delete = DELETE drop-statistics.

---

### Task 6: Partition Statistics CRUD — Resource

**Files:**
- Expand: `internal/client/statistics.go`
- Create: `internal/resources/partition_statistics/resource.go`
- Create: `internal/resources/partition_statistics/resource_test.go`

Same pattern as Task 5 but for partition-level statistics.

---

### Task 7: Lineage — Client + Data source

**Files:**
- Expand: `internal/client/lineage.go` (currently empty? create if not exists)
- Create: `internal/models/lineage.go`
- Create: `internal/datasources/lineage/` (if it doesn't exist)

**API Reference:** [Gravitino Lineage API](https://gravitino.apache.org/docs/next/api/rest/post-run-event)

The lineage endpoint is a POST for reporting run events — more of an action than a resource. A data source may not fit well. Consider just adding client methods for now.

---

### Task 8: Remaining small gaps

**Files to modify/create:**

1. **Catalog test connection:** `client/catalog.go` — add `TestCatalogConnection(metalake, catalog, body) (*TestConnectionResponse, error)`
   - POST `/metalakes/:m/catalogs/:c/test`
   - Could be a data source: `gravitino_catalog_test`

2. **Fileset files:** `client/fileset.go` — add `ListFilesetFiles(metalake, catalog, schema, fileset) (*FilesetFileListResponse, error)`
   - GET `/metalakes/:m/catalogs/:c/schemas/:s/filesets/:fs/files`
   - Could be a data source: `gravitino_fileset_files`

3. **SetMetalakeInUse:** `client/metalake.go` — add `SetMetalakeInUse(metalake, inUse bool) error`

4. **AssociatePoliciesForObject:** `client/policy.go` — add `AssociatePoliciesForObject(metalake, objType, objName string, policyNames []string) error`

5. **SetPolicy enable/disable:** `client/policy.go` — add `SetPolicy(metalake, policyName string, enabled bool) error`

6. **ListMetadataObjectsForPolicy:** `client/policy.go` — add `ListMetadataObjectsForPolicy(metalake, policyName string) (*MetadataObjectsListResponse, error)`

---

## Task dependency graph

```
Task 1 (model versions client) → Task 2 (model versions resource)
Task 3 (job templates client)  → Task 4 (job templates resource)
Task 5 (statistics)            → independent
Task 6 (partition statistics)  → independent
Task 7 (lineage)               → independent
Task 8 (small gaps)            → independent
```

Tasks 1+2, 3+4, 5, 6, 7, 8 are all parallel tracks.

---

## Validation checklist

After each task/track:

- [ ] `go build ./...` passes
- [ ] `go test -v -cover ./internal/...` passes
- [ ] All new resources have `client.NewResourceError` error handling
- [ ] All new resources have `client.IsNotFoundError` 404 checks
- [ ] All new resources have `tflog.Debug` at CRUD boundaries
- [ ] All enum fields have `stringvalidator.OneOf`
- [ ] Import with dot-separated ID parsing + valid/invalid tests
- [ ] Resources registered in `provider.go`
- [ ] `make generate` run for docs
- [ ] `golangci-lint run --config .github/golangci.yml ./...` passes
