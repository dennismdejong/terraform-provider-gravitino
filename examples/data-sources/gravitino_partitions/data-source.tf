data "gravitino_partitions" "example" {
  metalake = "example_metalake"
  catalog  = "hive_catalog"
  schema   = "example_schema"
  table    = "orders"
}
