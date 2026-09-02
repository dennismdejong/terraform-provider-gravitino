resource "gravitino_policy" "data_retention" {
  metalake = gravitino_metalake.example.name
  name     = "data_retention"
  comment  = "Retain data for compliance"

  policy_type = "custom"
  enabled     = true

  supported_object_types = ["CATALOG", "SCHEMA", "TABLE"]

  properties = {
    key1 = "value1"
  }

  custom_rules = {
    rule1 = "123"
  }
}
