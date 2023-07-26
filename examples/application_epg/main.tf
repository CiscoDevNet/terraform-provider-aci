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

resource "aci_tenant" "terraform_tenant" {
  name        = "static_leaf_tenant"
  description = "This tenant is created by terraform"
}

resource "aci_application_profile" "test_ap" {
  tenant_dn   = aci_tenant.terraform_tenant.id
  name        = "test"
  description = "from terraform"
  name_alias  = "test_ap"
  prio        = "level1"
}

resource "aci_application_epg" "fooapplication_epg" {
  application_profile_dn = aci_application_profile.test_ap.id
  name                   = "demo_epg"
  description            = "from terraform"
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
  relation_fv_rs_node_att {
    node_dn              = "topology/pod-1/node-108"
    encap                = "vlan-100"
    description          = "this is desc for static leaf"
    deployment_immediacy = "lazy"
    mode                 = "regular"
  }
}