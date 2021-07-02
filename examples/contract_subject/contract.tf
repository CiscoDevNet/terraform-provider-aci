resource "aci_tenant" "tenant_for_contract" {
  name        = "tenant_for_contract"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_contract" "democontract" {
  tenant_dn                = aci_tenant.tenant_for_contract.id
  name                     = "test_tf_contract"
  description              = "This contract is created by terraform ACI provider"
  scope                    = "context"
  target_dscp              = "VA"
  prio                     = "unspecified"
  relation_vz_rs_graph_att = aci_rest.rest_abs_graph.id # Relation to vnsAbsGraph class. Cardinality - N_TO_ONE
}

resource "aci_contract_subject" "foocontract_subject" {
  contract_dn   = aci_contract.democontract.id
  description   = "from terraform"
  name          = "demo_subject"
  annotation    = "tag_subject"
  cons_match_t  = "AtleastOne"
  name_alias    = "alias_subject"
  prio          = "level1"
  prov_match_t  = "AtleastOne"
  rev_flt_ports = "yes"
  target_dscp   = "CS0"
}

