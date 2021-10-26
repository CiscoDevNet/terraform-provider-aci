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

resource "aci_tenant" "footenant" {
  description = "sample aci_tenant from terraform"
  name        = "demo_tenant"
  annotation  = "tag_tenant"
  name_alias  = "alias_tenant"
}

resource "aci_l3_outside" "fool3_outside" {
  tenant_dn      = aci_tenant.footenant.id
  description    = "aci_l3_outside"
  name           = "demo_l3out"
  annotation     = "tag_l3out"
  enforce_rtctrl = ["export", "import"]
  name_alias     = "alias_out"
  target_dscp    = "unspecified"
}

resource "aci_external_network_instance_profile" "fooexternal_network_instance_profile" {
  l3_outside_dn  = aci_l3_outside.fool3_outside.id
  description    = "aci_external_network_instance_profile"
  name           = "demo_inst_prof"
  annotation     = "tag_network_profile"
  exception_tag  = "2"
  flood_on_encap = "disabled"
  match_t        = "AtleastOne"
  name_alias     = "alias_profile"
  pref_gr_memb   = "exclude"
  prio           = "level1"
  target_dscp    = "exclude"
}

resource "aci_l3_ext_subnet" "foosubnet" {
  external_network_instance_profile_dn = aci_external_network_instance_profile.fooexternal_network_instance_profile.id
  description                          = "Sample L3 External subnet"
  ip                                   = "10.0.3.28/27"
  aggregate                            = "shared-rtctrl"
  annotation                           = "tag_ext_subnet"
  name_alias                           = "alias_ext_subnet"
  scope                                = ["import-security"]
}