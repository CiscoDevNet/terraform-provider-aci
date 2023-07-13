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

resource "aci_power_supply_redundancy_policy" "foo_ps_redundancy_policy" {
  name        = "example_ps_redudancy_policy"
  admin_rdn_m = "comb"
  annotation  = "example"
  name_alias  = "example"
}
