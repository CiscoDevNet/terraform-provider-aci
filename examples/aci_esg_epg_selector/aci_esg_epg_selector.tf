resource "aci_tenant" "terraform_tenant" {
  name        = "tf_tenant"
  description = "This tenant is created by terraform"
}

resource "aci_application_profile" "terraform_ap" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "tf_ap"
}

resource "aci_vrf" "vrf" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "vrf"
}

resource "aci_bridge_domain" "demobd" {
  tenant_dn                   = aci_tenant.terraform_tenant.id
  name                        = "test_tf_bd"
  description                 = "This bridge domain is created by terraform ACI provider"
  optimize_wan_bandwidth      = "no"
  arp_flood                   = "no"
  ep_clear                    = "no"
  relation_fv_rs_ctx          = aci_vrf.vrf.id
  intersite_bum_traffic_allow = "yes"
  intersite_l2_stretch        = "yes"
  ip_learning                 = "yes"
  limit_ip_learn_to_subnets   = "yes"
  mcast_allow                 = "yes"
  multi_dst_pkt_act           = "bd-flood"
  bridge_domain_type          = "regular"
  unicast_route               = "yes"
  unk_mac_ucast_act           = "flood"
  unk_mcast_act               = "flood"
  vmac                        = "not-applicable"
}

resource "aci_application_epg" "app2" {
  application_profile_dn = aci_application_profile.terraform_ap.id
  name                   = "tf_ap2"
  relation_fv_rs_bd      = aci_bridge_domain.demobd.id
}

# create ESG
resource "aci_endpoint_security_group" "terraform_esg" {
  application_profile_dn = aci_application_profile.terraform_ap.id
  name                   = "tf_esg"
}

# create ESG EPG selector
resource "aci_endpoint_security_group_epg_selector" "terraform_epg_selector" {
  endpoint_security_group_dn = aci_endpoint_security_group.terraform_esg.id
  match_epg_dn               = aci_application_epg.app2.id
  description                = "EPG selector"
}
