resource "aci_tenant" "tenant_for_contract" {
  name        = "tenant_for_contract"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_contract" "democontract" {
  tenant_dn   = aci_tenant.tenant_for_contract.id
  name        = "test_tf_contract"
  description = "This contract is created by terraform ACI provider"
  scope       = "context"
  target_dscp = "VA"
  prio        = "unspecified"
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

resource "aci_filter" "allow_https" {
  tenant_dn = aci_tenant.tenant_for_contract.id
  name      = "allow_https"
}

resource "aci_filter" "allow_icmp" {
  tenant_dn = aci_tenant.tenant_for_contract.id
  name      = "allow_icmp"
}

resource "aci_contract_subject" "Web_subject1" {
  contract_dn                  = aci_contract.democontract.id
  name                         = "Subject"
  relation_vz_rs_subj_filt_att = [aci_filter.allow_https.id, aci_filter.allow_icmp.id]
}

resource "aci_l4_l7_service_graph_template" "service_graph" {
  tenant_dn = aci_tenant.tenant_for_contract.id
  name      = "test_service_graph"
}

resource "aci_l4_l7_service_graph_template" "service_graph2" {
  tenant_dn = aci_tenant.tenant_for_contract.id
  name      = "tf_service_graph"
}

// apply_both_directions is selected [yes] by default and there is only one filter required
resource "aci_contract_subject" "contract_subject" {
  contract_dn   = aci_contract.democontract.id
  name          = "contract_subject"
  rev_flt_ports = "no"
}

resource "aci_filter" "test_filter" {
  tenant_dn   = aci_tenant.tenant_for_contract.id
  name        = "test_tf_filter"
  description = "This filter is created by terraform ACI provider."
}

resource "aci_filter" "tf_filter" {
  tenant_dn   = aci_tenant.tenant_for_contract.id
  name        = "tf_filter"
  description = "This filter is created by terraform ACI provider."
}

// apply_both_directions is not selected and there are two filters (consumer_to_provider and provider_to_consumer)
resource "aci_contract_subject" "contract_subject_2" {
  contract_dn           = aci_contract_subject.contract_subject.contract_dn
  name                  = "contract_subject_2"
  rev_flt_ports         = "no"
  apply_both_directions = "no"
  consumer_to_provider {
    prio                             = "unspecified"
    target_dscp                      = "AF41"
    relation_vz_rs_in_term_graph_att = aci_l4_l7_service_graph_template.service_graph.id
    relation_vz_rs_filt_att {
      action = "deny"
      directives = ["log", "no_stats"]
      priority_override = "level2"
      filter_dn = aci_filter.test_filter.id
    }
    relation_vz_rs_filt_att {
      action = "permit"
      directives = ["log"]
      priority_override = "default"
      filter_dn = aci_filter.tf_filter.id
    }
  }
  provider_to_consumer {
    prio                              = "unspecified"
    target_dscp                       = "AF42"
    relation_vz_rs_out_term_graph_att = aci_l4_l7_service_graph_template.service_graph.id
    relation_vz_rs_filt_att {
      action = "deny"
      directives = ["log"]
      priority_override = "level2"
      filter_dn = aci_filter.test_filter.id
    }
    relation_vz_rs_filt_att {
      action = "permit"
      directives = ["log", "no_stats"]
      priority_override = "default"
      filter_dn = aci_filter.tf_filter.id
    }
  }
}

resource "aci_contract_subject_one_way_filter" "contract_subject_filter" {
  contract_subject_dn = tolist(aci_contract_subject.contract_subject_2.consumer_to_provider)[0].id
  action = "permit"
  directives = ["log"]
  priority_override = "default"
  filter_dn = aci_filter.test_filter.id
}

resource "aci_contract_subject_one_way_filter" "contract_subject_filter2" {
  contract_subject_dn = tolist(aci_contract_subject.contract_subject_2.provider_to_consumer)[0].id
  action = "permit"
  directives = ["log"]
  priority_override = "default"
  filter_dn = aci_filter.test_filter.id
}

data "aci_contract_subject" "example" {
  contract_dn = aci_contract_subject.contract_subject_2.contract_dn
  name        = aci_contract_subject.contract_subject_2.name
}

output "name" {
  value = data.aci_contract_subject.example
}