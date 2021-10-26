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
resource "aci_tenant" "example" {
  name = "test_acc_tenant"
}

resource "aci_l3_outside" "example" {
  tenant_dn   = aci_tenant.example.id
  name        = "demo_l3out"
  target_dscp = "CS0"
}

resource "aci_external_network_instance_profile" "fooexternal_network_instance_profile" {
  l3_outside_dn  = aci_l3_outside.example.id
  description    = "%s"
  name           = "demo_inst_prof"
  annotation     = "tag_network_profile"
  exception_tag  = "2"
  flood_on_encap = "disabled"
  match_t        = "ALL"
  name_alias     = "alias_profile"
  pref_gr_memb   = "exclude"
  prio           = "level1"
  target_dscp    = "unspecified"
}