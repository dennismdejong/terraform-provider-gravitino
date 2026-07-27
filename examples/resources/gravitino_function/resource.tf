resource "gravitino_function" "parse_json" {
  metalake      = gravitino_metalake.example.name
  catalog       = gravitino_catalog.hive.name
  schema        = gravitino_schema.example.name
  name          = "parse_json"
  comment       = "UDF for parsing JSON strings"
  function_body = "CREATE FUNCTION parse_json AS 'com.example.udf.JsonParser' USING JAR 's3://jars/udf-library.jar'"
}

resource "gravitino_function" "mask_pii" {
  metalake      = gravitino_metalake.example.name
  catalog       = gravitino_catalog.hive.name
  schema        = gravitino_schema.example.name
  name          = "mask_pii"
  comment       = "Mask PII fields with configurable masking character"
  function_body = "CREATE FUNCTION mask_pii AS 'com.example.udf.PiiMasker' USING JAR 's3://jars/udf-library.jar'"
  properties = {
    "type"          = "udf"
    "return_type"   = "string"
    "deterministic" = "true"
  }
}
