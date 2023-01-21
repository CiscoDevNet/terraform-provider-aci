resource "aci_tenant" "tenant_for_contract" {
  name        = "tenant_for_contract"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_l4_l7_service_graph_template" "rest_abs_graph" {
  tenant_dn = aci_tenant.tenant_for_contract.id
  name      = "testgraph"
}

// Creating a contract
resource "aci_contract" "web_contract" {
  tenant_dn   = aci_tenant.tenant_for_contract.id
  name        = "test_tf_contract"
  description = "This contract is created by terraform ACI provider"
  scope       = "context"
  target_dscp = "VA"
  prio        = "unspecified"
}

// Creating contract subject to connect contract to filters and filter entries in filter.tf file
resource "aci_contract_subject" "web_subject" {
  contract_dn                   = aci_contract.web_contract.id
  name                          = "Subject"
  relation_vz_rs_subj_graph_att = aci_l4_l7_service_graph_template.rest_abs_graph.id
  relation_vz_rs_subj_filt_att  = [aci_filter.allow_https.id, aci_filter.allow_icmp.id]
}

// Creating a contract with a filter and filter entry
resource "aci_contract" "complex_contract" {
  tenant_dn   = aci_tenant.tenant_for_contract.id
  name        = "test_tf_complex_contract"
  description = "This contract is created by terraform ACI provider"
  scope       = "context"
  target_dscp = "VA"
  prio        = "unspecified"
  filter {
    filter_name = "complex_contract_filter"
    filter_entry {
      filter_entry_name = "complex_contract_filter_entry"
      description = "My complex entry description from Terraform"
    }
  }
}