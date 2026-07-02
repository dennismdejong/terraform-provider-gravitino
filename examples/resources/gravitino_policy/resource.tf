resource "gravitino_policy" "data_retention" {
  metalake      = gravitino_metalake.example.name
  resource_type = "tables"
  resource      = "orders"
  name          = "data_retention"
  effect        = "allow"
}

resource "gravitino_policy" "read_only" {
  metalake      = gravitino_metalake.example.name
  resource_type = "schemas"
  resource      = "*"
  name          = "read_only"
  effect        = "deny"
  actions       = ["write", "delete"]
  subjects      = ["guest-user"]
  condition     = "context.role != 'admin'"
}
