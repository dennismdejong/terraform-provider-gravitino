resource "gravitino_table" "json_cols" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.hive.name
  schema   = gravitino_schema.example.name
  name     = "users"
  comment  = "User table with JSON column definitions"
  columns = [
    { "name" = "id", "type" = "integer", "nullable" = false },
    { "name" = "name", "type" = "string", "nullable" = true },
    { "name" = "email", "type" = "string" },
  ]
}

resource "gravitino_table" "partitioned" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.hive.name
  schema   = gravitino_schema.example.name
  name     = "orders"
  comment  = "Partitioned orders table"
  columns = [
    { "name" = "order_id", "type" = "string" },
    { "name" = "customer_id", "type" = "integer" },
    { "name" = "amount", "type" = "double" },
  ]
  partition_strategy = "HASH"
  partition_buckets  = 16
}
