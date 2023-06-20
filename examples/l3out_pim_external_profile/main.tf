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

resource "aci_tenant" "terraform_tenant" {
  name = "terraform_tenant"
}

resource "aci_l3_outside" "foo_l3_outside" {
  tenant_dn      = aci_tenant.terraform_tenant.id
  name           = "foo_l3_outside"
  enforce_rtctrl = ["export", "import"]
  target_dscp    = "unspecified"
  mpls_enabled   = "yes"
}

resource "aci_pim_external_profile" "example" {
  l3_outside_dn  = aci_l3_outside.foo_l3_outside.id
  enabled_af = ["ipv4-mcast"]
}
