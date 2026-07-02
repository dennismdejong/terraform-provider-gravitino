resource "gravitino_fileset" "managed" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.fileset.name
  schema   = gravitino_schema.example.name
  name     = "ml_datasets"
  type     = "managed"
  comment  = "Managed fileset for ML training data"
}

resource "gravitino_fileset" "external" {
  metalake        = gravitino_metalake.example.name
  catalog         = gravitino_catalog.fileset.name
  schema          = gravitino_schema.example.name
  name            = "raw_data"
  type            = "external"
  storage_location = "s3a://datalake/raw"
  comment         = "External fileset for raw data ingestion"
  properties = {
    "read-only" = "true"
  }
}
