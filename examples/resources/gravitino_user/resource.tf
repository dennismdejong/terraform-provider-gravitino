# Minimal user
resource "gravitino_user" "analyst" {
  metalake = gravitino_metalake.example.name
  name     = "data_analyst"
}

# User with roles assigned
resource "gravitino_user" "engineer" {
  metalake = gravitino_metalake.example.name
  name     = "data_engineer"
  roles    = ["admin", "readonly"]
}

# User with custom properties
resource "gravitino_user" "viewer" {
  metalake = gravitino_metalake.example.name
  name     = "report_viewer"
  roles    = ["readonly"]
}
