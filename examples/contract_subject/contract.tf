resource "aci_tenant" "tenant_for_contract" {
  name        = "tenant_for_contract"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_contract" "democontract" {
  tenant_dn   = aci_tenant.tenant_for_contract.id
  name        = "test_tf_contract"
  description = "This contract is created by terraform ACI provider"
}

resource "aci_filter" "allow_https" {
  tenant_dn = aci_tenant.tenant_for_contract.id
  name      = "allow_https"
}

resource "aci_filter" "allow_icmp" {
  tenant_dn = aci_tenant.tenant_for_contract.id
  name      = "allow_icmp"
}

// If apply_both_directions is "yes" [by default], the filters defined in the subject are applied for both directions.
resource "aci_contract_subject" "contract_subject" {
  contract_dn                  = aci_contract.democontract.id
  description                  = "from terraform"
  name                         = "demo_subject"
  annotation                   = "tag_subject"
  cons_match_t                 = "AtleastOne"
  prio                         = "level1"
  prov_match_t                 = "AtleastOne"
  rev_flt_ports                = "yes"
  target_dscp                  = "CS0"
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

// If apply_both_directions is "no", the filters for each direction (consumer-to-provider and provider-to-consumer) are defined independently.
resource "aci_contract_subject" "contract_subject_2" {
  contract_dn           = aci_contract.democontract.id
  name                  = "contract_subject_2"
  rev_flt_ports         = "no"
  apply_both_directions = "no"
  consumer_to_provider {
    prio                             = "unspecified"
    target_dscp                      = "AF41"
    relation_vz_rs_in_term_graph_att = aci_l4_l7_service_graph_template.service_graph.id
    relation_vz_rs_filt_att {
      action            = "deny"
      directives        = ["log", "no_stats"]
      priority_override = "level2"
      filter_dn         = aci_filter.allow_https.id
    }
    relation_vz_rs_filt_att {
      action            = "permit"
      directives        = ["log"]
      priority_override = "default"
      filter_dn         = aci_filter.allow_icmp.id
    }
  }
  provider_to_consumer {
    prio                              = "unspecified"
    target_dscp                       = "AF42"
    relation_vz_rs_out_term_graph_att = aci_l4_l7_service_graph_template.service_graph.id
    relation_vz_rs_filt_att {
      action            = "deny"
      directives        = ["log"]
      priority_override = "level2"
      filter_dn         = aci_filter.allow_https.id
    }
    relation_vz_rs_filt_att {
      action            = "permit"
      directives        = ["log", "no_stats"]
      priority_override = "default"
      filter_dn         = aci_filter.allow_icmp.id
    }
  }
}

data "aci_contract_subject" "example" {
  contract_dn = aci_contract_subject.contract_subject_2.contract_dn
  name        = aci_contract_subject.contract_subject_2.name
}

output "name" {
  value = data.aci_contract_subject.example
}