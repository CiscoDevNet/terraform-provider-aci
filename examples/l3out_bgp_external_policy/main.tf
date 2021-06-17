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

resource "aci_l3out_bgp_external_policy" "example" {

  l3_outside_dn = aci_l3_outside.example.id
  annotation    = "example"
  name_alias    = "example"

}
