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

resource "aci_vrf_snmp_context" "example" {
  vrf_dn     = aci_vrf.example.id
  name       = "example"
  annotation = "orchestrator:terraform"
  name_alias = "example_name_alias"
}