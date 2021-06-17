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

resource "aci_bgp_best_path_policy" "foobgp_best_path_policy" {
  tenant_dn  = aci_tenant.example.id
  name       = "example"
  annotation = "example"
  ctrl       = "asPathMultipathRelax"
  name_alias = "example"
}
