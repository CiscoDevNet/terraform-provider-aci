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

resource "aci_tacacs_source" "example" {
  parent_dn   = "uni/fabric/moncommon"
  name        = "example"
  annotation  = "orchestrator:terraform"
  incl        = ["audit", "session"]
  min_sev     = "info"
  name_alias  = "tacacs_source_alias"
  description = "From Terraform"
}