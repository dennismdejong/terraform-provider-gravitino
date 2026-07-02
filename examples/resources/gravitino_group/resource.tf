# Minimal group
resource "gravitino_group" "analysts" {
  metalake = gravitino_metalake.example.name
  name     = "analytics_team"
}

# Group with roles
resource "gravitino_group" "engineers" {
  metalake = gravitino_metalake.example.name
  name     = "engineering"
  roles    = ["admin"]
}

# Read-only group
resource "gravitino_group" "readers" {
  metalake = gravitino_metalake.example.name
  name     = "readonly_users"
  roles    = ["readonly"]
}
