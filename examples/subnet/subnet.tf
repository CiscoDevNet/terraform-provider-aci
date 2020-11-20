resource "aci_tenant" "tenant_for_subnet" {
  name        = "tenant_for_subnet"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_bridge_domain" "bd_for_subnet" {
  tenant_dn   = "${aci_tenant.tenant_for_subnet.id}"
  name        = "bd_for_subnet"
  description = "This bridge domain is created by terraform ACI provider"
  mac         = "00:22:BD:F8:19:FF"
}

resource "aci_subnet" "demosubnet" {
  parent_dn                           = "${aci_bridge_domain.bd_for_subnet.id}"
  ip                                  = "10.0.3.28/27"
  scope                               = ["private"]
  description                         = "This subject is created by terraform"
  ctrl                                = ["unspecified"]
  preferred                           = "no"
  virtual                             = "yes"
  relation_fv_rs_bd_subnet_to_profile = "${aci_rest.rest_rt_ctrl_profile.id}" # Relation to rtctrlProfile class. Cardinality - N_TO_ONE.
  relation_fv_rs_bd_subnet_to_out     = ["${aci_rest.rest_l3_ext_out.id}"]    # Relation to l3extOut class. Cardinality - N_TO_M.
  relation_fv_rs_nd_pfx_pol           = "${aci_rest.rest_nd_pfx_pol.id}"      # Relation to ndPfxPol class. Cardinality - N_TO_ONE.
}
