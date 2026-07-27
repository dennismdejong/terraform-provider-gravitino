## 0.3.1 (2026-07-27)

BREAKING CHANGES: None

ENHANCEMENTS:
- **Documentation:** Provider docs now show all 7 auth method examples with HCL snippets
- **Resource examples:** Enriched with more realistic configurations
  - Table: added `sort_orders` and `distribution` example
  - View: added `view_def` (SQL) and properties
  - Function: added `function_body` and properties
  - Partition: multiple partitions across different tables
- **Complete deployment example:** New `examples/complete/main.tf` showing all resources working together in a realistic data platform scenario
- **AGENTS.md:** Updated with auth architecture, missing resources (model_version, job_template), and gokrb5 dependency

## 0.3.0 (2026-07-27)

BREAKING CHANGES: None

FEATURES:
- **Provider authentication:** Full support for all 4 Gravitino auth methods
  - `simple` — OS user or `GRAVITINO_USER` environment variable
  - `basic` — HTTP Basic authentication (unchanged)
  - `oauth` — Static bearer token AND client credentials flow with auto-refresh
  - `kerberos` — SPNEGO authentication via keytab or ticket cache

ENHANCEMENTS:
- Provider config expanded with 8 new attributes: `oauth_client_id`, `oauth_client_secret`, `oauth_server_uri`, `oauth_token_path`, `oauth_scope`, `kerberos_principal`, `kerberos_keytab`, `kerberos_use_ticket_cache`
- Modular auth architecture for easier extensibility
- OAuth2 token refresh with thread-safe caching (90% expiry margin)
- Kerberos SPNEGO with 401 challenge-response retry handling
- New `none` auth option for explicit no-authentication configuration

NOTES:
- `auth = "basic"` and `auth = "oauth"` with `oauth_token` remain fully backward compatible
- Empty or unset `auth` still works (no authentication)
- `gokrb5/v8` added as dependency for Kerberos support

## 0.1.0 (Unreleased)

BREAKING CHANGES: None

FEATURES:
- **New Resource:** `gravitino_metalake`
- **New Resource:** `gravitino_catalog`
- **New Resource:** `gravitino_schema`
- **New Resource:** `gravitino_table` with full column, sort order, distribution, partitioning, and index support
- **New Resource:** `gravitino_tag`
- **New Resource:** `gravitino_fileset`
- **New Resource:** `gravitino_topic`
- **New Resource:** `gravitino_view`
- **New Resource:** `gravitino_function`
- **New Resource:** `gravitino_model`
- **New Resource:** `gravitino_partition`
- **New Resource:** `gravitino_policy`
- **New Resource:** `gravitino_job`
- Provider supports OAuth2 and HTTP Basic authentication
- All resources support import
- Enum validators for constrained string attributes
