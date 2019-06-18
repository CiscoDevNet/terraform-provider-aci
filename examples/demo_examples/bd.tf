resource "aci_tenant" "tenant_for_bridge_domain" {
  name        = "tenant_for_bd"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_bridge_domain" "demobd" {
  tenant_dn   = "${aci_tenant.tenant_for_bridge_domain.id}"
  name        = "test_tf_bd"
  description = "This bridge domain is created by terraform ACI provider"
  mac         = "00:22:BD:F8:19:FF"
}
