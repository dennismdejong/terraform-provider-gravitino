terraform {
  required_providers {
    gravitino = {
      source = "gravitino/gravitino"
    }
  }
}

provider "gravitino" {
  uri      = "http://localhost:8090"
  auth     = "basic"
  username = "admin"
  password = "admin"
}
