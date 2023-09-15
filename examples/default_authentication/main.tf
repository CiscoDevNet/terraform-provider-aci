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

resource "aci_default_authentication" "example" {
  annotation     = "orchestrator:terraform"
  fallback_check = "false"
  realm          = "local"
  realm_sub_type = "default"
  name_alias     = "example_name_alias"
  description    = "from terraform"
}