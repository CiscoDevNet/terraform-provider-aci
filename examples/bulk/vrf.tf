resource "aci_vrf" "test_tf_vrf-1" {
  tenant_dn              = aci_tenant.tenant_for_benchmark.id
  name                   = "test_tf_vrf-1"
  description            = "from terraform"
  annotation             = "tag_vrf"
  bd_enforced_enable     = "no"
  ip_data_plane_learning = "enabled"
  knw_mcast_act          = "permit"
  name_alias             = "alias_vrf-1"
  pc_enf_dir             = "egress"
  pc_enf_pref            = "unenforced"
}
