resource "gravitino_tag" "pii" {
  metalake = gravitino_metalake.example.name
  name     = "PII"
  comment  = "Personally Identifiable Information"
}

resource "gravitino_tag" "confidential" {
  metalake   = gravitino_metalake.example.name
  name       = "confidential"
  comment    = "Confidential data requiring special access"
  properties = {
    "classification" = "internal"
    "retention"      = "90d"
  }
}
