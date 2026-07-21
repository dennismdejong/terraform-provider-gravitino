resource "gravitino_job_template" "etl" {
  metalake = gravitino_metalake.example.name
  name     = "daily_etl"
  template = "spark"
  parameters = {
    "main_class" = "com.example.ETLJob"
    "jar"        = "s3://jars/etl-job.jar"
    "args"       = "--date=${date}"
  }
  comment = "Daily ETL job template"
}
