resource "gravitino_job" "daily_etl" {
  metalake   = gravitino_metalake.example.name
  name       = "daily_etl"
  template   = "spark_etl"
  schedule   = "0 2 * * *"
  parameters = {
    "input_path"  = "/data/raw"
    "output_path" = "/data/processed"
  }
}
