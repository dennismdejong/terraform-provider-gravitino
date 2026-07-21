data "gravitino_model_version" "example" {
  metalake = "example_metalake"
  catalog  = "ml_catalog"
  schema   = "example_schema"
  model    = "fraud_detector"
  version  = "v1.0"
}
