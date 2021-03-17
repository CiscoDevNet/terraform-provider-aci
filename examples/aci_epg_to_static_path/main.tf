terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

# provider "aci" {
#   username    = ""
#   private_key = ""
#   cert_name   = ""
#   url         = ""
#   insecure    = true
# }
