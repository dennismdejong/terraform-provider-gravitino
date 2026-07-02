resource "gravitino_function" "example" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.hive.name
  schema   = gravitino_schema.example.name
  name     = "parse_json"
  comment  = "UDF for parsing JSON strings"
}
