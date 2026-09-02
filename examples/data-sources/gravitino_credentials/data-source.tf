data "gravitino_credentials" "example" {
  metalake      = "example_metalake"
  resource_type = "TABLE"
  resource      = "hive_catalog.example_schema.users"
}
