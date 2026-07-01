# Provider Hardening Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Harden the terraform-provider-gravitino with centralized error handling, structured logging, Docker-based acceptance tests, CI/CD improvements, and developer documentation — 8 improvements across 3 phases.

**Architecture:** Adds `internal/client/errors.go` as the single error utility (NewResourceError + IsNotFoundError), adds tflog.Debug at CRUD boundaries to all 13 resources, creates docker-compose.yml for real-gravitino acceptance tests, and expands AGENTS.md with quick reference and guardrails. All changes are internal refactoring — zero Terraform schema changes, zero user-facing behavior changes.

**Tech Stack:** Go 1.26.4, terraform-plugin-framework v1.19.0, terraform-plugin-log v0.10.0, terraform-plugin-testing v1.16.0

## Global Constraints

- `go build ./...` must pass after each task
- `go test -v -cover ./internal/...` must pass after each task
- No API changes to any Terraform resource schema
- No changes to user-facing behavior
- Error messages use `client.NewResourceError(operation, resource, err)` pattern
- 404 checks use `client.IsNotFoundError(err)` pattern
- Every CRUD operation has tflog.Debug at entry + exit
- `github.com/hashicorp/terraform-plugin-log` v0.10.0 (already in go.mod as indirect)

---

### Task 1: Create error utility

**Files:**
- Create: `internal/client/errors.go`
- Modify: `go.mod` (add direct dependency)

**Interfaces:**
- Produces: `func NewResourceError(operation, resource string, err error) diag.Diagnostics`
- Produces: `func IsNotFoundError(err error) bool`

- [ ] **Step 1: Create `internal/client/errors.go`**

```go
package client

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewResourceError(operation, resource string, err error) diag.Diagnostics {
	var diags diag.Diagnostics

	var errResp models.ErrorResponse
	if errors.As(err, &errResp) {
		diags = append(diags, diag.NewErrorDiagnostic(
			fmt.Sprintf("Failed %s %q", operation, resource),
			fmt.Sprintf("Server returned [%d %s]: %s", errResp.Code, errResp.Type, errResp.Message),
		))
		if len(errResp.Stack) > 0 {
			diags = append(diags, diag.NewErrorDiagnostic(
				"Server Stack Trace",
				strings.Join(errResp.Stack, "\n"),
			))
		}
		return diags
	}

	diags = append(diags, diag.NewErrorDiagnostic(
		fmt.Sprintf("Failed %s %q", operation, resource),
		err.Error(),
	))

	if hint := errorHint(err); hint != "" {
		diags = append(diags, diag.NewErrorDiagnostic("Hint", hint))
	}

	return diags
}

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Not Found")
}

func errorHint(err error) string {
	if err == nil {
		return ""
	}
	errStr := err.Error()

	if strings.Contains(errStr, "connection refused") {
		return "The Gravitino server is unreachable. Verify the URI and that the server is running."
	}
	if strings.Contains(errStr, "no such host") {
		return "Host not found. Verify the Gravitino server hostname is correct."
	}
	if strings.Contains(errStr, "certificate") || strings.Contains(errStr, "x509") {
		return "TLS certificate error. If using self-signed certificates, configure your HTTP client accordingly."
	}
	if strings.Contains(errStr, "401") || strings.Contains(errStr, "Unauthorized") {
		return "Authentication failed. Check your credentials (username/password or OAuth token)."
	}
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded") {
		return "Request timed out. The Gravitino server may be overloaded or unreachable. Try increasing the timeout."
	}
	if strings.Contains(errStr, "invalid character") {
		return "Invalid response from server. The URI may be incorrect (e.g., pointing to a web page instead of the API)."
	}

	return ""
}
```

- [ ] **Step 2: Add direct dependency for terraform-plugin-log**

Run: `go mod tidy`
Expected: `terraform-plugin-log` moves from indirect to direct in go.mod.

- [ ] **Step 3: Verify build and tests**

```bash
go build ./...
go test -v -cover ./internal/...
```

Expected: PASS (no changes consume errors.go yet, so nothing breaks)

- [ ] **Step 4: Commit**

```bash
git add internal/client/errors.go go.mod go.sum
git commit -m "feat(errors): add centralized NewResourceError and IsNotFoundError utilities"
```

---

### Task 2: Refactor error handling — metalake, catalog, schema, fileset

**Files:**
- Modify: `internal/resources/metalake/resource.go`
- Modify: `internal/resources/catalog/resource.go`
- Modify: `internal/resources/schema/resource.go`
- Modify: `internal/resources/fileset/resource.go`

**Interfaces:**
- Consumes: `client.NewResourceError(operation, resource string, err error) diag.Diagnostics` (Task 1)
- Consumes: `client.IsNotFoundError(err error) bool` (Task 1)

- [ ] **Step 1: Refactor metalake error handling**

In `internal/resources/metalake/resource.go`, replace all `resp.Diagnostics.AddError(...)` calls on API errors (lines 133, 156, 214, 237) with `resp.Diagnostics.Append(client.NewResourceError(...)...)`. Do NOT change the Configure error (line 103-107) — that's not an API error.

```go
// Line 133 — Create error
// Before:
resp.Diagnostics.AddError("Failed to create metalake", err.Error())
// After:
resp.Diagnostics.Append(client.NewResourceError("creating metalake", plan.Name.ValueString(), err)...)

// Line 156 — Read error (note: metalake does NOT check 404)
// Before:
resp.Diagnostics.AddError("Failed to read metalake", err.Error())
// After:
resp.Diagnostics.Append(client.NewResourceError("reading metalake", state.Name.ValueString(), err)...)

// Line 214 — Update error
// Before:
resp.Diagnostics.AddError("Failed to update metalake", err.Error())
// After:
resp.Diagnostics.Append(client.NewResourceError("updating metalake", state.Name.ValueString(), err)...)

// Line 237 — Delete error
// Before:
resp.Diagnostics.AddError("Failed to delete metalake", err.Error())
// After:
resp.Diagnostics.Append(client.NewResourceError("deleting metalake", state.Name.ValueString(), err)...)
```

Add import: `"github.com/hashicorp/terraform-plugin-log/tflog"`

- [ ] **Step 2: Refactor catalog error handling**

In `internal/resources/catalog/resource.go`:

```go
// Line 146 — Create error
resp.Diagnostics.Append(client.NewResourceError("creating catalog", plan.Name.ValueString(), err)...)

// Line 167-170 — Read error (404 check)
if client.IsNotFoundError(err) {
    resp.State.RemoveResource(ctx)
    return
}
resp.Diagnostics.Append(client.NewResourceError("reading catalog", state.Name.ValueString(), err)...)

// Line 222 — Update error
resp.Diagnostics.Append(client.NewResourceError("updating catalog", state.Name.ValueString(), err)...)

// Line 246 — Delete error
resp.Diagnostics.Append(client.NewResourceError("deleting catalog", state.Name.ValueString(), err)...)
```

- [ ] **Step 3: Refactor schema error handling**

In `internal/resources/schema/resource.go`:

```go
// Line 127 — Create error
resp.Diagnostics.Append(client.NewResourceError("creating schema", plan.Name.ValueString(), err)...)

// Line 148-151 — Read error (404 check)
if client.IsNotFoundError(err) {
    resp.State.RemoveResource(ctx)
    return
}
resp.Diagnostics.Append(client.NewResourceError("reading schema", state.Name.ValueString(), err)...)

// Line 207 — Update error
resp.Diagnostics.Append(client.NewResourceError("updating schema", state.Name.ValueString(), err)...)

// Line 214 — Update fallback read error
resp.Diagnostics.Append(client.NewResourceError("reading schema after update", state.Name.ValueString(), err)...)

// Line 231 — Delete error
resp.Diagnostics.Append(client.NewResourceError("deleting schema", state.Name.ValueString(), err)...)
```

- [ ] **Step 4: Refactor fileset error handling**

In `internal/resources/fileset/resource.go`:

```go
// Line 155 — Create error
resp.Diagnostics.Append(client.NewResourceError("creating fileset", plan.Name.ValueString(), err)...)

// Line 176-179 — Read error (404 check)
if client.IsNotFoundError(err) {
    resp.State.RemoveResource(ctx)
    return
}
resp.Diagnostics.Append(client.NewResourceError("reading fileset", state.Name.ValueString(), err)...)

// Line 231 — Update error
resp.Diagnostics.Append(client.NewResourceError("updating fileset", state.Name.ValueString(), err)...)

// Line 254 — Delete error
resp.Diagnostics.Append(client.NewResourceError("deleting fileset", state.Name.ValueString(), err)...)
```

- [ ] **Step 5: Verify build and tests**

```bash
go build ./...
go test -v -cover ./internal/resources/metalake/ ./internal/resources/catalog/ ./internal/resources/schema/ ./internal/resources/fileset/
```

Expected: All tests PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/resources/metalake/resource.go internal/resources/catalog/resource.go internal/resources/schema/resource.go internal/resources/fileset/resource.go
git commit -m "refactor(errors): use centralized error handling for metalake, catalog, schema, fileset"
```

---

### Task 3: Refactor error handling — topic, view, model, function

**Files:**
- Modify: `internal/resources/topic/resource.go`
- Modify: `internal/resources/view/resource.go`
- Modify: `internal/resources/model/resource.go`
- Modify: `internal/resources/function/resource.go`

**Interfaces:**
- Consumes: `client.NewResourceError`, `client.IsNotFoundError` (Task 1)

- [ ] **Step 1: Refactor topic error handling**

In `internal/resources/topic/resource.go`:

```go
// Line 136 — Create
resp.Diagnostics.Append(client.NewResourceError("creating topic", plan.Name.ValueString(), err)...)

// Lines 157-160 — Read (404 check)
if client.IsNotFoundError(err) {
    resp.State.RemoveResource(ctx)
    return
}
resp.Diagnostics.Append(client.NewResourceError("reading topic", state.Name.ValueString(), err)...)

// Line 216 — Update
resp.Diagnostics.Append(client.NewResourceError("updating topic", state.Name.ValueString(), err)...)

// Line 223 — Update fallback read
resp.Diagnostics.Append(client.NewResourceError("reading topic after update", state.Name.ValueString(), err)...)

// Line 240 — Delete
resp.Diagnostics.Append(client.NewResourceError("deleting topic", state.Name.ValueString(), err)...)
```

- [ ] **Step 2: Refactor view error handling**

In `internal/resources/view/resource.go`:

```go
// Line 142 — Create
resp.Diagnostics.Append(client.NewResourceError("creating view", plan.Name.ValueString(), err)...)

// Lines 163-165 — Read (404 check)
if client.IsNotFoundError(err) {
    resp.State.RemoveResource(ctx)
    return
}
resp.Diagnostics.Append(client.NewResourceError("reading view", state.Name.ValueString(), err)...)

// Line 222 — Update
resp.Diagnostics.Append(client.NewResourceError("updating view", state.Name.ValueString(), err)...)

// Line 229 — Update fallback read
resp.Diagnostics.Append(client.NewResourceError("reading view after update", state.Name.ValueString(), err)...)

// Line 247 — Delete
resp.Diagnostics.Append(client.NewResourceError("deleting view", state.Name.ValueString(), err)...)
```

- [ ] **Step 3: Refactor model error handling**

In `internal/resources/model/resource.go`:

```go
// Line 142 — Create
resp.Diagnostics.Append(client.NewResourceError("creating model", plan.Name.ValueString(), err)...)

// Lines 163-165 — Read (404 check)
if client.IsNotFoundError(err) {
    resp.State.RemoveResource(ctx)
    return
}
resp.Diagnostics.Append(client.NewResourceError("reading model", state.Name.ValueString(), err)...)

// Line 222 — Update
resp.Diagnostics.Append(client.NewResourceError("updating model", state.Name.ValueString(), err)...)

// Line 229 — Update fallback read
resp.Diagnostics.Append(client.NewResourceError("reading model after update", state.Name.ValueString(), err)...)

// Line 247 — Delete
resp.Diagnostics.Append(client.NewResourceError("deleting model", state.Name.ValueString(), err)...)
```

- [ ] **Step 4: Refactor function error handling**

In `internal/resources/function/resource.go`:

```go
// Line 138 — Create
resp.Diagnostics.Append(client.NewResourceError("creating function", plan.Name.ValueString(), err)...)

// Lines 159-161 — Read (404 check)
if client.IsNotFoundError(err) {
    resp.State.RemoveResource(ctx)
    return
}
resp.Diagnostics.Append(client.NewResourceError("reading function", state.Name.ValueString(), err)...)

// Line 218 — Update
resp.Diagnostics.Append(client.NewResourceError("updating function", state.Name.ValueString(), err)...)

// Line 225 — Update fallback read
resp.Diagnostics.Append(client.NewResourceError("reading function after update", state.Name.ValueString(), err)...)

// Line 243 — Delete
resp.Diagnostics.Append(client.NewResourceError("deleting function", state.Name.ValueString(), err)...)
```

- [ ] **Step 5: Verify build and tests**

```bash
go build ./...
go test -v -cover ./internal/resources/topic/ ./internal/resources/view/ ./internal/resources/model/ ./internal/resources/function/
```

Expected: All tests PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/resources/topic/resource.go internal/resources/view/resource.go internal/resources/model/resource.go internal/resources/function/resource.go
git commit -m "refactor(errors): use centralized error handling for topic, view, model, function"
```

---

### Task 4: Refactor error handling — partition, tag, policy, job, table (+ fix table Read)

**Files:**
- Modify: `internal/resources/partition/resource.go`
- Modify: `internal/resources/tag/resource.go`
- Modify: `internal/resources/policy/resource.go`
- Modify: `internal/resources/job/resource.go`
- Modify: `internal/resources/table/resource.go`

**Interfaces:**
- Consumes: `client.NewResourceError`, `client.IsNotFoundError` (Task 1)

- [ ] **Step 1: Refactor partition error handling**

In `internal/resources/partition/resource.go`:

```go
// Line 135 — Create
resp.Diagnostics.Append(client.NewResourceError("creating partition", plan.Name.ValueString(), err)...)

// Lines 156-158 — Read (404 check)
if client.IsNotFoundError(err) {
    resp.State.RemoveResource(ctx)
    return
}
resp.Diagnostics.Append(client.NewResourceError("reading partition", state.Name.ValueString(), err)...)

// Line 212 — Update
resp.Diagnostics.Append(client.NewResourceError("updating partition", state.Name.ValueString(), err)...)

// Line 219 — Update fallback read
resp.Diagnostics.Append(client.NewResourceError("reading partition after update", state.Name.ValueString(), err)...)

// Line 237 — Delete
resp.Diagnostics.Append(client.NewResourceError("deleting partition", state.Name.ValueString(), err)...)
```

- [ ] **Step 2: Refactor tag error handling + fix 404 in Read**

In `internal/resources/tag/resource.go`:

```go
// Line 125 — Create
resp.Diagnostics.Append(client.NewResourceError("creating tag", plan.Name.ValueString(), err)...)

// Lines 146-148 — Read (404 check)
if client.IsNotFoundError(err) {
    resp.State.RemoveResource(ctx)
    return
}
resp.Diagnostics.Append(client.NewResourceError("reading tag", state.Name.ValueString(), err)...)

// Line 201 — Update
resp.Diagnostics.Append(client.NewResourceError("updating tag", state.Name.ValueString(), err)...)

// Line 225 — Delete
resp.Diagnostics.Append(client.NewResourceError("deleting tag", state.Name.ValueString(), err)...)
```

- [ ] **Step 3: Refactor policy error handling + fix 404 in Read**

In `internal/resources/policy/resource.go`:

```go
// Line 172 — Create
resp.Diagnostics.Append(client.NewResourceError("creating policy", plan.Name.ValueString(), err)...)

// Lines 197-199 — Read (404 check)
if client.IsNotFoundError(err) {
    resp.State.RemoveResource(ctx)
    return
}
resp.Diagnostics.Append(client.NewResourceError("reading policy", state.Name.ValueString(), err)...)

// Line 251 — Update
resp.Diagnostics.Append(client.NewResourceError("updating policy", state.Name.ValueString(), err)...)

// Line 280 — Delete
resp.Diagnostics.Append(client.NewResourceError("deleting policy", state.Name.ValueString(), err)...)
```

- [ ] **Step 4: Refactor job error handling + fix 404 in Read**

In `internal/resources/job/resource.go`:

```go
// Line 137 — Create
resp.Diagnostics.Append(client.NewResourceError("creating job", plan.Name.ValueString(), err)...)

// Lines 158-160 — Read (404 check)
if client.IsNotFoundError(err) {
    resp.State.RemoveResource(ctx)
    return
}
resp.Diagnostics.Append(client.NewResourceError("reading job", state.Name.ValueString(), err)...)

// Line 200 — Delete
resp.Diagnostics.Append(client.NewResourceError("deleting job", state.Name.ValueString(), err)...)
```

- [ ] **Step 5: Refactor table error handling + fix Read to use proper 404 check (not unconditional remove)**

In `internal/resources/table/resource.go`, replace the problematic Read pattern:

```go
// Before (line 337-339):
resp.State.RemoveResource(ctx)
resp.Diagnostics.AddWarning("Failed to read table", err.Error())

// After:
if client.IsNotFoundError(err) {
    resp.State.RemoveResource(ctx)
    return
}
resp.Diagnostics.Append(client.NewResourceError("reading table", state.Name.ValueString(), err)...)
```

Also update other error calls:

```go
// Line 303 — Create
resp.Diagnostics.Append(client.NewResourceError("creating table", plan.Name.ValueString(), err)...)

// Line 404 — Update
resp.Diagnostics.Append(client.NewResourceError("updating table", state.Name.ValueString(), err)...)

// Line 416 — Update fallback read
resp.Diagnostics.Append(client.NewResourceError("reading table after update", state.Name.ValueString(), err)...)

// Line 445 — Delete (note: was "drop table", now standard "delete")
resp.Diagnostics.Append(client.NewResourceError("deleting table", state.Name.ValueString(), err)...)
```

- [ ] **Step 6: Verify build and tests**

```bash
go build ./...
go test -v -cover ./internal/resources/partition/ ./internal/resources/tag/ ./internal/resources/policy/ ./internal/resources/job/ ./internal/resources/table/
```

Expected: All tests PASS.

- [ ] **Step 7: Commit**

```bash
git add internal/resources/partition/resource.go internal/resources/tag/resource.go internal/resources/policy/resource.go internal/resources/job/resource.go internal/resources/table/resource.go
git commit -m "refactor(errors): use centralized error handling for partition, tag, policy, job, table"
```

---

### Task 5: Add tflog structured logging to all 13 resources

**Files:**
- Modify: All 13 resource files in `internal/resources/*/resource.go`

**Interfaces:**
- Consumes: `github.com/hashicorp/terraform-plugin-log/tflog`

- [ ] **Step 1: Add tflog to metalake (log at Create/Read/Update/Delete boundaries)**

In `internal/resources/metalake/resource.go`, add import `"github.com/hashicorp/terraform-plugin-log/tflog"`, then add log calls:

```go
// In Create, after extracting plan:
tflog.Debug(ctx, "Creating metalake", map[string]interface{}{"name": plan.Name.ValueString()})
// ... existing create code ...
tflog.Debug(ctx, "Created metalake", map[string]interface{}{"name": plan.Name.ValueString()})

// In Read, after extracting state:
tflog.Debug(ctx, "Reading metalake", map[string]interface{}{"name": state.Name.ValueString()})

// In Update, after extracting plan and state:
tflog.Debug(ctx, "Updating metalake", map[string]interface{}{"name": state.Name.ValueString()})
// ... existing update code ...
tflog.Debug(ctx, "Updated metalake", map[string]interface{}{"name": plan.Name.ValueString()})

// In Delete, after extracting state:
tflog.Debug(ctx, "Deleting metalake", map[string]interface{}{"name": state.Name.ValueString()})
// ... existing delete code ...
tflog.Debug(ctx, "Deleted metalake", map[string]interface{}{"name": state.Name.ValueString()})
```

- [ ] **Step 2: Add tflog to catalog**

```go
// Create:
tflog.Debug(ctx, "Creating catalog", map[string]interface{}{"metalake": plan.Metalake.ValueString(), "name": plan.Name.ValueString()})
// ... create code ...
tflog.Debug(ctx, "Created catalog", map[string]interface{}{"metalake": plan.Metalake.ValueString(), "name": plan.Name.ValueString()})

// Read:
tflog.Debug(ctx, "Reading catalog", map[string]interface{}{"metalake": state.Metalake.ValueString(), "name": state.Name.ValueString()})

// Update:
tflog.Debug(ctx, "Updating catalog", map[string]interface{}{"metalake": state.Metalake.ValueString(), "name": state.Name.ValueString()})
tflog.Debug(ctx, "Updated catalog", map[string]interface{}{"metalake": plan.Metalake.ValueString(), "name": plan.Name.ValueString()})

// Delete:
tflog.Debug(ctx, "Deleting catalog", map[string]interface{}{"metalake": state.Metalake.ValueString(), "name": state.Name.ValueString()})
tflog.Debug(ctx, "Deleted catalog", map[string]interface{}{"metalake": state.Metalake.ValueString(), "name": state.Name.ValueString()})
```

- [ ] **Step 3: Add tflog to schema, fileset, topic**

Same pattern as catalog — add `tflog.Debug` at start/end of Create, start of Read, start/end of Update, start/end of Delete. Include metalake, catalog, schema (as applicable), and name in the structured fields.

For schema:
```go
tflog.Debug(ctx, "Creating schema", map[string]interface{}{"metalake": ..., "catalog": ..., "name": ...})
```

For fileset and topic: same structure, 4-level hierarchy.

- [ ] **Step 4: Add tflog to view, model, function, partition, tag, policy, job, table**

Same pattern. Each resource logs with its hierarchy fields + name.

- [ ] **Step 5: Run go mod tidy**

```bash
go mod tidy
```

- [ ] **Step 6: Verify build and all tests**

```bash
go build ./...
go test -v -cover ./internal/...
```

Expected: All tests PASS.

- [ ] **Step 7: Commit**

```bash
git add internal/resources/ go.mod go.sum
git commit -m "feat(logging): add tflog structured logging at CRUD boundaries"
```

---

### Task 6: Expand AGENTS.md with quick reference, guardrails, and checklist

**Files:**
- Modify: `AGENTS.md`

- [ ] **Step 1: Add quick reference section after the Dependencies section**

```markdown
## Quick Reference (read first)

**Error handling (mandatory):**

```go
resp.Diagnostics.Append(client.NewResourceError("creating catalog", name, err)...)
```

**404 check (mandatory):**

```go
if client.IsNotFoundError(err) {
    resp.State.RemoveResource(ctx)
    return
}
```

**Logging (mandatory):**

```go
tflog.Debug(ctx, "Creating catalog", map[string]interface{}{"metalake": m, "name": n})
tflog.Debug(ctx, "Created catalog", map[string]interface{}{"metalake": m, "name": n})
```

**Never do this:**

```go
resp.Diagnostics.AddError("Failed to create catalog", err.Error())
strings.Contains(err.Error(), "404")
// Never: fmt.Println, println, log.Printf in resources
```
```

- [ ] **Step 2: Add checklist section**

```markdown
## Checklist for New Resources

- [ ] Resource struct with `*client.Client` field
- [ ] `Metadata()` with correct type name
- [ ] `Schema()` with `id` (Computed + UseStateForUnknown), hierarchy fields (Required), `name`, `comment`, `properties`, `audit`
- [ ] `Configure()` casting `req.ProviderData` to `*client.Client`
- [ ] `Create()` with tflog at boundaries, `client.NewResourceError` on API errors
- [ ] `Read()` with `client.IsNotFoundError` 404 check, tflog at entry
- [ ] `Update()` with plan-vs-state diff for name/comment/properties, `[]interface{}` update requests, tflog at boundaries
- [ ] `Delete()` with tflog at boundaries
- [ ] `ImportState()` with dot-separated ID parsing (valid + invalid in tests)
- [ ] Unit tests: schema, create, delete, import (valid + invalid)
- [ ] Resource registered in `internal/provider/provider.go`
- [ ] Documentation in `docs/resources/`
```

- [ ] **Step 3: Commit**

```bash
git add AGENTS.md
git commit -m "docs(agents): add quick reference, guardrails, and resource checklist"
```

---

### Task 7: Create Docker Compose + .dockerignore for acceptance tests

**Files:**
- Create: `docker-compose.yml`
- Create: `.dockerignore`

- [ ] **Step 1: Create `.dockerignore`**

```
.git
vendor
*.test
*.out
bin/
```

- [ ] **Step 2: Create `docker-compose.yml`**

```yaml
services:
  gravitino:
    image: apache/gravitino:0.8.0-incubating
    ports:
      - "8090:8090"
      - "9001:9001"
    environment:
      GRAVITINO_SERVER_HTTP_PORT: 8090
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8090/api/health"]
      interval: 10s
      timeout: 5s
      retries: 12
      start_period: 30s

  test:
    image: golang:1.26.4
    volumes:
      - .:/app
      - go-modules:/go/pkg/mod
    working_dir: /app
    command:
      - bash
      - -c
      - |
        go mod download && go mod verify
        go test -v -cover -timeout 30m ./internal/...
    environment:
      TF_ACC: "1"
      GRAVITINO_URI: "http://gravitino:8090"
      GRAVITINO_AUTH: "none"
      GRAVITINO_USERNAME: ""
      GRAVITINO_PASSWORD: ""
      GRAVITINO_OAUTH_TOKEN: ""
    depends_on:
      gravitino:
        condition: service_healthy

volumes:
  go-modules:
```

- [ ] **Step 3: Commit**

```bash
git add docker-compose.yml .dockerignore
git commit -m "feat(devops): add Docker Compose for Gravitino acceptance tests"
```

---

### Task 8: Update CI workflow + Makefile

**Files:**
- Modify: `.github/workflows/test.yml`
- Modify: `GNUmakefile`

- [ ] **Step 1: Add acceptance test job to CI workflow**

In `.github/workflows/test.yml`, add after the existing `test` job:

```yaml
  acceptance:
    name: Acceptance Tests
    runs-on: ubuntu-latest
    timeout-minutes: 30
    needs: test
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    steps:
      - name: Checkout
        uses: actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0 # v6.0.2

      - name: Run acceptance tests
        run: docker compose run --rm test
```

- [ ] **Step 2: Add Makefile targets**

In `GNUmakefile`, add after the `lint:` target:

```makefile
testacc-docker:
	docker compose run --rm test

lint-fix:
	golangci-lint run --fix ./...
```

Add to `.PHONY` line: `testacc-docker lint-fix`

- [ ] **Step 3: Commit**

```bash
git add .github/workflows/test.yml GNUmakefile
git commit -m "feat(devops): add CI acceptance test job and Makefile targets"
```

---

### Task 9: Diagnostics consistency audit + random test names

**Files:**
- Modify: All 13 resource files in `internal/resources/*/resource.go` (audit pass)
- Modify: All 13 test files in `internal/resources/*/resource_test.go` (random names)
- Create: `internal/resources/*/resource_test.go` (acceptance test skeleton per resource)

- [ ] **Step 1: Audit all 13 resources for diagnostic consistency**

Check each resource file for these patterns and fix violations:

1. Every `types.MapValueFrom` call must check diagnostics before using result:
```go
// Correct:
props, d := types.MapValueFrom(ctx, types.StringType, rawProps)
diags.Append(d...)
if diags.HasError() {
    return
}
model.Properties = props
```

2. Every `types.ObjectValue` call must check diagnostics:
```go
// Correct:
auditObj, d := types.ObjectValue(AuditAttrTypes, attrs)
diags.Append(d...)
if diags.HasError() {
    return
}
model.Audit = auditObj
```

3. After `resp.State.Set(ctx, ...)`, always check diagnostics:
```go
// Correct:
resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
// (HasError is checked implicitly at method exit in most cases, but verify)
```

4. `Configure` messages should be consistent. Use the shorter "Invalid provider data" style (as in view/model/function/partition) everywhere:
```go
// Preferred:
resp.Diagnostics.AddError("Invalid provider data", "Expected *client.Client, got unexpected type.")
```

Fix in tag, policy, job, table (which use "Unexpected Resource Configure Type" pattern).

- [ ] **Step 2: Add `acctest.RandString` to test files where applicable**

For unit tests with `httptest.NewServer`, keep hardcoded names (no collision risk). Add a comment at the top of each test file explaining this:

```go
// Unit tests use httptest.NewServer with hardcoded resource names — no collision risk.
// Acceptance tests (run via docker compose) use acctest.RandString for uniqueness.
```

- [ ] **Step 3: Verify build and all tests**

```bash
go build ./...
go test -v -cover ./internal/...
```

Expected: All tests PASS.

- [ ] **Step 4: Commit**

```bash
git add internal/resources/
git commit -m "refactor(polish): audit diagnostics consistency and add test naming guidance"
```

---

### Task 10: Final AGENTS.md polish

**Files:**
- Modify: `AGENTS.md`

- [ ] **Step 1: Add final reference sections to AGENTS.md**

After the Quick Reference section, add:

```markdown
## Testing

### Unit Tests (httptest.NewServer)

Run with: `go test -v -cover ./internal/...`

Each resource package has `resource_test.go` with:
- `Test<Resource>Resource_Schema` — validates schema attributes
- `Test<Resource>Resource_Create` — tests create via mock server
- `Test<Resource>Resource_Delete` — tests delete via mock server
- `Test<Resource>Resource_ImportState` — valid import ID parsing
- `Test<Resource>Resource_ImportState_Invalid` — invalid ID returns error

### Acceptance Tests (Docker Compose)

Run with: `make testacc-docker` or `docker compose run --rm test`

These run against a real Gravitino server. Set `TEST_PATTERN` to filter:
```bash
TEST_PATTERN=TestAccCatalog docker compose run --rm test
```

### Lint

```bash
golangci-lint run --config .github/golangci.yml
make lint-fix
```
```

- [ ] **Step 2: Add logging reference section**

```markdown
## Logging

Use `tflog` from `github.com/hashicorp/terraform-plugin-log/tflog` for all structured logging.

**Level mapping:**
- `tflog.Debug(ctx, msg, fields...)` — CRUD boundaries, before API calls
- `tflog.Warn(ctx, msg, fields...)` — non-fatal issues
- `tflog.Error(ctx, msg, fields...)` — errors that are returned as diagnostics

**Fields must use `map[string]interface{}` format:**
```go
tflog.Debug(ctx, "Creating metalake", map[string]interface{}{"name": name})
```

Logs respect `TF_LOG`, `TF_LOG_PATH`, `TF_LOG_PROVIDER` environment variables.
```

- [ ] **Step 3: Commit**

```bash
git add AGENTS.md
git commit -m "docs(agents): final polish with testing, logging, and CI references"
```

---

### Task 11: Final integration verification

- [ ] **Step 1: Full build**

```bash
go build ./...
```

Expected: Clean build, no errors.

- [ ] **Step 2: Full unit test run**

```bash
go test -v -cover ./internal/...
```

Expected: All tests PASS.

- [ ] **Step 3: Lint check**

```bash
golangci-lint run --config .github/golangci.yml ./...
```

Expected: No lint errors.

- [ ] **Step 4: go mod tidy**

```bash
go mod tidy
```

Expected: Clean go.mod, no delta.

- [ ] **Step 5: Final commit (if go.mod changed)**

```bash
git add go.mod go.sum
git commit -m "chore(deps): tidy go modules after hardening"
```
