data "gravitino_partition" "example" {
  metalake = "example_metalake"
  catalog  = "hive_catalog"
  schema   = "example_schema"
  table    = "orders"
  name     = "2024_q1"
}
