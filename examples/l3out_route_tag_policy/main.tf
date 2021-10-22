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

resource "aci_l3out_route_tag_policy" "example" {

  tenant_dn  = aci_tenant.example.id
  name  = "example"
  annotation  = "example"
  name_alias  = "example"
  tag  = "1"
  description = "from terraform"

}
