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

resource "aci_bd_dhcp_label" "foo_bd_dhcp_label" {

  bridge_domain_dn = aci_bridge_domain.foo_bridge_domain.id
  name             = "example_bd_dhcp_label"
  annotation       = "example"
  name_alias       = "example"
  owner            = "tenant"
  tag              = "aqua"
}