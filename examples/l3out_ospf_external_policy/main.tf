terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

#configure provider with your cisco aci credentials.
provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <cloud APIC URL>
  insecure = true
}

resource "aci_l3out_ospf_external_policy" "example" {

  l3_outside_dn     = aci_l3_outside.example.id
  annotation        = "example"
  area_cost         = "1"
  area_ctrl         = ["redistribute", "summary"]
  area_id           = "0.0.0.1"
  area_type         = "nssa"
  multipod_internal = "no"
  name_alias        = "example"

}

