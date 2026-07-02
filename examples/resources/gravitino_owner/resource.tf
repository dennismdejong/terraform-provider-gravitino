# Set a user as owner of a catalog
resource "gravitino_owner" "catalog_owner" {
  metalake         = gravitino_metalake.example.name
  object_type      = "CATALOG"
  object_full_name = gravitino_catalog.hive.name
  owner_name       = "data_engineer"
  owner_type       = "USER"
}

# Set a group as owner of a schema
resource "gravitino_owner" "schema_owner" {
  metalake         = gravitino_metalake.example.name
  object_type      = "SCHEMA"
  object_full_name = "${gravitino_catalog.hive.name}.${gravitino_schema.example.name}"
  owner_name       = "engineering"
  owner_type       = "GROUP"
}
