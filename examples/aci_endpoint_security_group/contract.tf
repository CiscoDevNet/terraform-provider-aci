resource "aci_contract" "rs_prov_contract" {
  tenant_dn   = aci_tenant.terraform_tenant.id
  name        = "rs_prov_contract"
  description = "This contract is created by terraform ACI provider"
  scope       = "tenant"
  target_dscp = "VA"
  prio        = "unspecified"
}

resource "aci_contract" "rs_cons_contract" {
  tenant_dn   = aci_tenant.terraform_tenant.id
  name        = "rs_cons_contract"
  description = "This contract is created by terraform ACI provider"
  scope       = "tenant"
  target_dscp = "VA"
  prio        = "unspecified"
}

resource "aci_contract" "intra_epg_contract" {
  tenant_dn   = aci_tenant.terraform_tenant.id
  name        = "intra_epg_contract"
  description = "This contract is created by terraform ACI provider"
  scope       = "tenant"
  target_dscp = "VA"
  prio        = "unspecified"
}

resource "aci_contract" "exported_contract" {
  tenant_dn   = aci_tenant.tenant_export_cons.id
  name        = "exported_contract"
  description = "This contract is exported to tenant_export_cons"
}