# List all metalakes
data "gravitino_metalakes" "all" {}

# Get a specific metalake
data "gravitino_metalake" "example" {
  name = "example_metalake"
}

# List catalogs in a metalake
data "gravitino_catalogs" "hive" {
  metalake = "example_metalake"
}

# List schemas in a catalog
data "gravitino_schemas" "default" {
  metalake = "example_metalake"
  catalog  = "hive_catalog"
}

# List tables in a schema
data "gravitino_tables" "all" {
  metalake = "example_metalake"
  catalog  = "hive_catalog"
  schema   = "default"
}

# Get a specific table
data "gravitino_table" "users" {
  metalake = "example_metalake"
  catalog  = "hive_catalog"
  schema   = "default"
  name     = "users"
}

# Health check
data "gravitino_health" "status" {}
data "gravitino_liveness" "alive" {}
data "gravitino_readiness" "ready" {}
