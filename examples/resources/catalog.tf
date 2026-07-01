resource "gravitino_catalog" "hive" {
  metalake = gravitino_metalake.example.name
  name     = "hive_catalog"
  type     = "relational"
  provider = "hive"
  comment  = "Hive catalog for data warehouse"
  properties = {
    "metastore.uris" = "thrift://localhost:9083"
  }
}

resource "gravitino_catalog" "fileset" {
  metalake = gravitino_metalake.example.name
  name     = "fileset_catalog"
  type     = "fileset"
  comment  = "Fileset catalog for data lake"
}
