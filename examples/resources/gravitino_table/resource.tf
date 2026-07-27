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
  comment  = "Orders table with hash partitioning"
  columns = [
    { "name" = "order_id", "type" = "string" },
    { "name" = "customer_id", "type" = "integer" },
    { "name" = "amount", "type" = "double" },
  ]
  partition_strategy = "HASH"
  partition_buckets  = 16
}

resource "gravitino_table" "sorted" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.hive.name
  schema   = gravitino_schema.example.name
  name     = "events"
  comment  = "Events table with sort order and distribution"
  columns = [
    { "name" = "event_id",    "type" = "string" },
    { "name" = "event_type",  "type" = "string" },
    { "name" = "user_id",     "type" = "integer" },
    { "name" = "event_time",  "type" = "timestamp" },
  ]
  partition_strategy = "RANGE"
  partition_buckets  = 12
  sort_orders = [
    { "name" = "event_time", "direction" = "DESC" },
    { "name" = "event_type", "direction" = "ASC" },
  ]
  distribution = {
    strategy = "HASH"
    buckets  = 4
    partition_keys = [
      { "name" = "user_id", "type" = "integer" },
    ]
  }
}
