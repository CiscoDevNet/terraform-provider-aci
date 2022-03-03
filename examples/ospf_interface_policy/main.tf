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
  ctrl         = ["advert-subnet", "bfd"]
  dead_intvl   = "40"
  hello_intvl  = "10"
  name_alias   = "alias_ospf"
  nw_t         = "unspecified"
  pfx_suppress = "inherit"
  prio         = "1"
  rexmit_intvl = "5"
  xmit_delay   = "1"
}
