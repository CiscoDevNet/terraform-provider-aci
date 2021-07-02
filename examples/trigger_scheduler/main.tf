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

resource "aci_trigger_scheduler" "example" {
  name        = "example"
  annotation  = "example"
  description = "from terraform"
  name_alias  = "example"
}
