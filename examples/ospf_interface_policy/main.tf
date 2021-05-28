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
resource "aci_tenant" "dev_tenant" {
  name        = "tf_test_rel_tenant"
  description = "This tenant is created by terraform"
}
resource "aci_ospf_interface_policy" "fooospf_interface_policy" {
  tenant_dn    = aci_tenant.dev_tenant.id
  description  = "From Terraform"
  name         = "demo_ospfpol"
  annotation   = "tag_ospf"
  cost         = "unspecified"
  ctrl         = "unspecified"
  dead_intvl   = "40"
  hello_intvl  = "10"
  name_alias   = "alias_ospf"
  nw_t         = "unspecified"
  pfx_suppress = "inherit"
  prio         = "1"
  rexmit_intvl = "5"
  xmit_delay   = "1"
}
