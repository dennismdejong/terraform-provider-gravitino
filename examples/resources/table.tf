resource "gravitino_schema" "default" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.hive.name
  name     = "default"
  comment  = "Default schema"
}

resource "gravitino_table" "users" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.hive.name
  schema   = gravitino_schema.default.name
  name     = "users"
  comment  = "User table"

  column {
    name     = "id"
    type     = "integer"
    nullable = false
    auto_increment = true
  }
  column {
    name    = "username"
    type    = "varchar"
    length  = 255
    comment = "Unique username"
  }
  column {
    name   = "email"
    type   = "varchar"
    length = 512
  }
  column {
    name     = "created_at"
    type     = "timestamp"
    nullable = false
  }

  properties = {
    "format" = "parquet"
  }

  sort_order {
    field_name    = ["id"]
    direction     = "asc"
    null_ordering = "nulls_first"
  }

  distribution {
    strategy  = "hash"
    number    = 16
    func_args = ["id"]
  }

  partitioning {
    strategy   = "day"
    field_name = ["created_at"]
  }

  index {
    index_type  = "primary_key"
    name        = "pk_users"
    field_names = [["id"]]
  }
}
