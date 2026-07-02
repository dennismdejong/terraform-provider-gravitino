resource "gravitino_model" "example" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.ml.name
  schema   = gravitino_schema.example.name
  name     = "fraud_detector"
  comment  = "ML model for fraud detection"
}
