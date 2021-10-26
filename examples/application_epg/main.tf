
provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
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