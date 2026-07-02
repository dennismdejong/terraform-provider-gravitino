data "gravitino_function" "example" {
  metalake = "example_metalake"
  catalog  = "hive_catalog"
  schema   = "example_schema"
  name     = "parse_json"
}
