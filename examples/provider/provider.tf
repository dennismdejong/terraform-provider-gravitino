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

# Simple authentication (uses OS user or GRAVITINO_USER env)
provider "gravitino" {
  uri  = "http://localhost:8090"
  auth = "simple"
}

# OAuth2 static bearer token
provider "gravitino" {
  uri         = "http://localhost:8090"
  auth        = "oauth"
  oauth_token = var.gravitino_token
}

# OAuth2 client credentials flow
provider "gravitino" {
  uri                = "http://localhost:8090"
  auth               = "oauth"
  oauth_client_id     = var.oauth_client_id
  oauth_client_secret = var.oauth_client_secret
  oauth_server_uri    = "http://localhost:8177"
  oauth_token_path    = "/oauth2/token"
  oauth_scope         = "test"
}

# Kerberos authentication with keytab
provider "gravitino" {
  uri                = "http://gravitino.example.com"
  auth               = "kerberos"
  kerberos_principal  = "HTTP/gravitino.example.com@EXAMPLE.COM"
  kerberos_keytab     = "/etc/security/gravitino.keytab"
}

# Kerberos authentication with ticket cache
provider "gravitino" {
  uri                      = "http://gravitino.example.com"
  auth                     = "kerberos"
  kerberos_principal        = "HTTP/gravitino.example.com@EXAMPLE.COM"
  kerberos_use_ticket_cache = true
}
