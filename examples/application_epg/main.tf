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
  name = "demo_tenant"
}

resource "aci_application_profile" "test_ap" {
  tenant_dn   = aci_tenant.dev_tenant.id
  name        = "demo_ap"
  annotation  = "tag"
  description = "from terraform"
  name_alias  = "test_ap"
  prio        = "level1"
}

resource "aci_application_epg" "fooapplication_epg" {
  application_profile_dn = aci_application_profile.test_ap.id
  name                   = "demo_epg"
  description            = "from terraform"
  annotation             = "tag_epg"
  exception_tag          = "0"
  flood_on_encap         = "disabled"
  fwd_ctrl               = "none"
  has_mcast_source       = "no"
  is_attr_based_epg      = "no"
  match_t                = "AtleastOne"
  name_alias             = "alias_epg"
  pc_enf_pref            = "unenforced"
  pref_gr_memb           = "exclude"
  prio                   = "unspecified"
  shutdown               = "no"
}

/*
The following depicts an example to create and associate an application EPG with the common Tenant's BD and VRF
*/

data "aci_tenant" "common_tenant" {
  name = "common"
}

data "aci_vrf" "default_vrf" {
  tenant_dn = data.aci_tenant.common_tenant.id
  name      = "default"
}

resource "aci_bridge_domain" "test_bd" {
  tenant_dn          = data.aci_tenant.common_tenant.id
  name               = "common_test_bd"
  relation_fv_rs_ctx = data.aci_vrf.default_vrf.id
}

resource "aci_application_epg" "test_epg_common" {
  application_profile_dn = aci_application_profile.test_ap.id
  name                   = "common_test_epg"
  relation_fv_rs_bd      = aci_bridge_domain.test_bd.id
}
