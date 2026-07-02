resource "gravitino_topic" "events" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.kafka.name
  schema   = gravitino_schema.example.name
  name     = "user_events"
}

resource "gravitino_topic" "configured" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.kafka.name
  schema   = gravitino_schema.example.name
  name     = "audit_logs"
  comment  = "Topic for audit log events"
  properties = {
    "retention.ms" = "604800000"
    "partitions"   = "6"
  }
}
