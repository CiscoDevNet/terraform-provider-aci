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

resource "aci_l3out_bgp_protocol_profile" "example" {
  logical_node_profile_dn = aci_logical_node_profile.example.id
  annotation              = "example"
  name_alias              = "example"
}
