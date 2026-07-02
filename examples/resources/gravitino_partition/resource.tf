resource "gravitino_partition" "example" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.hive.name
  schema   = gravitino_schema.example.name
  table    = gravitino_table.example.name
  name     = "2024_q1"
}
