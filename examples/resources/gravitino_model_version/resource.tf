resource "gravitino_model_version" "v1" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.ml.name
  schema   = gravitino_schema.example.name
  model    = gravitino_model.example.name
  version  = "v1.0"
  uri      = "s3://models/fraud_detector/v1.0"
  aliases  = ["production", "latest"]
  comment  = "Initial production version"
}

resource "gravitino_model_version" "v2" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.ml.name
  schema   = gravitino_schema.example.name
  model    = gravitino_model.example.name
  version  = "v2.0"
  uri      = "s3://models/fraud_detector/v2.0"
  aliases  = ["staging"]
  comment  = "Improved model with new features"
  properties = {
    "accuracy" = "0.95"
    "framework" = "pytorch"
  }
}
