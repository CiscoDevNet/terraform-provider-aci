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

resource "aci_firmware_group" "example" {
  name  = "example"
  annotation  = "example"
  description = "from terraform"
  name_alias  = "example"
  firmware_group_type  = "range"
}