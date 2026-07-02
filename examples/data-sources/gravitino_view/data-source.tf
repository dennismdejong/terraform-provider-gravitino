data "gravitino_view" "example" {
  metalake = "example_metalake"
  catalog  = "hive_catalog"
  schema   = "example_schema"
  name     = "user_summary"
}
