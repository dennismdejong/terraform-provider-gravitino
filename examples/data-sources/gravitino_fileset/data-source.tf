data "gravitino_fileset" "example" {
  metalake = "example_metalake"
  catalog  = "hive_catalog"
  schema   = "example_schema"
  name     = "ml_datasets"
}
