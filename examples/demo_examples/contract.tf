resource "aci_tenant" "tenant_for_contract" {
  name        = "tenant_for_contract"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_contract" "democontract" {
  tenant_dn   = aci_tenant.tenant_for_contract.id
  name        = "test_tf_contract"
  description = "This contract is created by terraform ACI provider"
  scope       = "context"
  dscp        = "VA"
}
