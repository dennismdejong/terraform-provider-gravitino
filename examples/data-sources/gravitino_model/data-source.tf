data "gravitino_model" "example" {
  metalake = "example_metalake"
  catalog  = "ml_catalog"
  schema   = "example_schema"
  name     = "fraud_detector"
}
