resource "gravitino_view" "user_summary" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.hive.name
  schema   = gravitino_schema.example.name
  name     = "user_summary"
  comment  = "View aggregating user data"
  view_def = "SELECT u.id, u.name, COUNT(o.order_id) as order_count FROM users u LEFT JOIN orders o ON u.id = o.customer_id GROUP BY u.id, u.name"
}

resource "gravitino_view" "active_customers" {
  metalake = gravitino_metalake.example.name
  catalog  = gravitino_catalog.hive.name
  schema   = gravitino_schema.example.name
  name     = "active_customers"
  comment  = "Customers with orders in last 90 days"
  view_def = "SELECT DISTINCT u.id, u.name, u.email FROM users u INNER JOIN orders o ON u.id = o.customer_id WHERE o.order_date >= current_date - INTERVAL '90' DAY"
  properties = {
    "refresh" = "daily"
    "owner"   = "analytics-team"
  }
}
