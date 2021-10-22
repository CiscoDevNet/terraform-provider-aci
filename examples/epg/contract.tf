resource "aci_tenant" "tenant_for_contract" {
  name        = "tenant_for_contract"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_l4_l7_service_graph_template" "rest_abs_graph" {
  tenant_dn = aci_tenant.tenant_for_contract.id
  name      = "testgraph"
}

resource "aci_contract" "rs_prov_contract" {
  tenant_dn                = aci_tenant.tenant_for_contract.id
  name                     = "rs_prov_contract"
  description              = "This contract is created by terraform ACI provider"
  scope                    = "context"
  target_dscp              = "VA"
  prio                     = "unspecified"
  relation_vz_rs_graph_att = aci_l4_l7_service_graph_template.rest_abs_graph.id
}

resource "aci_contract" "rs_cons_contract" {
  tenant_dn                = aci_tenant.tenant_for_contract.id
  name                     = "rs_cons_contract"
  description              = "This contract is created by terraform ACI provider"
  scope                    = "context"
  target_dscp              = "VA"
  prio                     = "unspecified"
  relation_vz_rs_graph_att = aci_l4_l7_service_graph_template.rest_abs_graph.id
}

resource "aci_contract" "intra_epg_contract" {
  tenant_dn                = aci_tenant.tenant_for_contract.id
  name                     = "intra_epg_contract"
  description              = "This contract is created by terraform ACI provider"
  scope                    = "context"
  target_dscp              = "VA"
  prio                     = "unspecified"
  relation_vz_rs_graph_att = aci_l4_l7_service_graph_template.rest_abs_graph.id
}

// Taboo Contract
resource "aci_taboo_contract" "rest_taboo_con" {
  tenant_dn = aci_tenant.tenant_for_epg.id
  name      = "testcon"
}

// Imported Contract
resource "aci_imported_contract" "rest_vz_cons_if" {
  tenant_dn = aci_tenant.tenant_for_epg.id
  name      = "testcontract"
}