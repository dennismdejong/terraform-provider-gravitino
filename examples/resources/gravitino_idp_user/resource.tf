# Built-in IDP user for local authentication
resource "gravitino_idp_user" "alice" {
  name     = "alice"
  password = "Passw0rd-Alice12"
  enabled  = true
}

# Disabled user
resource "gravitino_idp_user" "legacy" {
  name     = "legacy_operator"
  password = "Passw0rd-Legacy12"
  enabled  = false
}
