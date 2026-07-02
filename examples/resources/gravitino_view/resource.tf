resource "gravitino_view" "example" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.hive.name
  schema   = gravitino_schema.example.name
  name     = "user_summary"
  comment  = "View aggregating user data"
}
