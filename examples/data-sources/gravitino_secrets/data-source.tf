data "gravitino_secrets" "example" {
  metalake      = "example_metalake"
  resource_type = "CATALOG"
  resource      = "hive_catalog"
}
