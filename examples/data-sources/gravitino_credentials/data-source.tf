data "gravitino_credentials" "example" {
  metalake      = "example_metalake"
  catalog       = "hive_catalog"
  schema        = "example_schema"
  resource_type = "TABLES"
  resource      = "users"
}
