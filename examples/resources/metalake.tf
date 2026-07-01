resource "gravitino_metalake" "example" {
  name    = "example_metalake"
  comment = "Example metalake for development"
  properties = {
    env = "dev"
  }
}
