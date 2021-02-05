
resource "aci_bridge_domain" "bd_for_rel" {
  tenant_dn   = aci_tenant.tenant_for_epg.id
  name        = "test_tf_bd_rel"
  description = "This bridge domain is created by terraform ACI provider"
  mac         = "00:22:BD:F8:19:FF"
}
