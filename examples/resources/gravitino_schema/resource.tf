resource "gravitino_schema" "example" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.hive.name
  name     = "example_schema"
}

resource "gravitino_schema" "full" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.hive.name
  name     = "analytics_schema"
  comment  = "Schema for analytics data"
  properties = {
    owner = "data-team"
  }
}
