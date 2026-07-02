data "gravitino_topic" "example" {
  metalake = "example_metalake"
  catalog  = "kafka_catalog"
  schema   = "example_schema"
  name     = "user_events"
}
