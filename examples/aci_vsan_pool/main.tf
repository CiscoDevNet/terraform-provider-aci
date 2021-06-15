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

resource "aci_vsan_pool" "example" {
  name        = "example"
  description = "from terraform"
  alloc_mode  = "static"
  annotation  = "example"
  name_alias  = "example"
}