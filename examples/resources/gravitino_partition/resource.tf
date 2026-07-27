resource "gravitino_partition" "q1_2024" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.hive.name
  schema   = gravitino_schema.example.name
  table    = gravitino_table.partitioned.name
  name     = "2024_q1"
}

resource "gravitino_partition" "q2_2024" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.hive.name
  schema   = gravitino_schema.example.name
  table    = gravitino_table.partitioned.name
  name     = "2024_q2"
}

resource "gravitino_partition" "eu_region" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.hive.name
  schema   = gravitino_schema.example.name
  table    = gravitino_table.sorted.name
  name     = "region_eu"
}

resource "gravitino_partition" "us_region" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.hive.name
  schema   = gravitino_schema.example.name
  table    = gravitino_table.sorted.name
  name     = "region_us"
}
