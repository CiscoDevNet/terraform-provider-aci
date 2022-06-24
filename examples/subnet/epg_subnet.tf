resource "aci_tenant" "foo_tenant" {
  name        = "foo_tenant"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_application_profile" "foo_app_profile" {
  tenant_dn   = aci_tenant.foo_tenant.id
  name        = "foo_app_profile"
  annotation  = "tag"
  description = "from terraform"
  name_alias  = "test_ap"
  prio        = "unspecified"
}

resource "aci_application_epg" "foo_epg" {
  application_profile_dn = aci_application_profile.foo_app_profile.id
  name                   = "foo_epg"
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

resource "aci_subnet" "foo_epg_subnet_next_hop_addr" {
  parent_dn     = aci_application_epg.foo_epg.id
  ip            = "10.0.3.29/32"
  scope         = ["private"]
  description   = "This subject is created by terraform"
  ctrl          = ["no-default-gateway"]
  preferred     = "no"
  virtual       = "yes"
  next_hop_addr = "10.0.3.30"
}

resource "aci_subnet" "foo_epg_subnet_anycast_mac" {
  parent_dn   = aci_application_epg.foo_epg.id
  ip          = "10.0.3.29/32"
  scope       = ["private"]
  description = "This subject is created by terraform"
  ctrl        = ["no-default-gateway"]
  preferred   = "no"
  virtual     = "yes"
  anycast_mac = "F0:1F:20:34:89:AB"
}

resource "aci_subnet" "foo_epg_subnet_msnlb_mcast_igmp" {
  parent_dn   = aci_application_epg.foo_epg.id
  ip          = "10.0.3.29/32"
  scope       = ["private"]
  description = "This subject is created by terraform"
  ctrl        = ["no-default-gateway"]
  preferred   = "no"
  virtual     = "yes"
  msnlb {
    mode  = "mode-mcast-igmp"
    group = "224.0.0.1" # Valid Multicast Address are 224.0.0.0 through 239.255.255.255
  }
}

resource "aci_subnet" "foo_epg_subnet_msnlb_mcast_static" {
  parent_dn   = aci_application_epg.foo_epg.id
  ip          = "10.0.3.29/32"
  scope       = ["private"]
  description = "This subject is created by terraform"
  ctrl        = ["no-default-gateway"]
  preferred   = "no"
  virtual     = "yes"
  msnlb {
    mode = "mode-mcast--static"
    mac  = "03:1F:20:34:89:AA"
  }
}

resource "aci_subnet" "foo_epg_subnet_msnlb_mode_uc" {
  parent_dn   = aci_application_epg.foo_epg.id
  ip          = "10.0.3.29/32"
  scope       = ["private"]
  description = "This subject is created by terraform"
  ctrl        = ["no-default-gateway"]
  preferred   = "no"
  virtual     = "yes"
  msnlb {
    mode = "mode-uc"
    mac  = "00:1F:20:34:89:AA"
  }
}
