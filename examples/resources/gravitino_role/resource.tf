# Minimal role (just a name, no privileges)
resource "gravitino_role" "readonly" {
  metalake = gravitino_metalake.example.name
  name     = "readonly"
}

# Role with securable objects and specific privileges
resource "gravitino_role" "data_engineer" {
  metalake = gravitino_metalake.example.name
  name     = "data_engineer"
  securable_objects = [
    {
      full_name = "*"
      type      = "CATALOG"
      privileges = [
        { name = "USE_CATALOG", condition = "ALLOW" }
      ]
    },
    {
      full_name = "analytics_schema"
      type      = "SCHEMA"
      privileges = [
        { name = "USE_SCHEMA",   condition = "ALLOW" },
        { name = "CREATE_TABLE", condition = "ALLOW" },
        { name = "SELECT_TABLE", condition = "ALLOW" },
      ]
    }
  ]
}

# Role with properties
resource "gravitino_role" "admin" {
  metalake = gravitino_metalake.example.name
  name     = "admin"
  properties = {
    "description" = "Full access role"
    "managed_by"  = "security-team"
  }
  securable_objects = [
    {
      full_name = "*"
      type      = "METALAKE"
      privileges = [
        { name = "CREATE_CATALOG", condition = "ALLOW" },
        { name = "MANAGE_USERS",   condition = "ALLOW" },
        { name = "MANAGE_GROUPS",  condition = "ALLOW" },
        { name = "CREATE_ROLE",    condition = "ALLOW" },
        { name = "MANAGE_GRANTS",  condition = "ALLOW" },
      ]
    }
  ]
}
