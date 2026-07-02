data "gravitino_roles" "example" {
  metalake      = "example_metalake"
  resource_type = "TABLES"
  resource      = "users"
}
