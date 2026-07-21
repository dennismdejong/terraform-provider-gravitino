data "gravitino_model_versions" "all" {
  metalake = "example_metalake"
  catalog  = "ml_catalog"
  schema   = "example_schema"
  model    = "fraud_detector"
}
