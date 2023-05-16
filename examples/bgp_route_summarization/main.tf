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

resource "aci_tenant" "tf_tenant" {
  name = "tf_tenant"
}

resource "aci_bgp_route_summarization" "bgp_rt_summ_pol" {
  tenant_dn             = aci_tenant.tf_tenant.id
  name                  = "bgp_rt_summ_pol"
  description           = "from terraform"
  attrmap               = "sample attrmap"
  ctrl                  = ["summary-only", "as-set"]
  address_type_controls = ["af-ucast", "af-mcast", "af-label-ucast"]
}
