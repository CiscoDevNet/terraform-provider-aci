terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = "admin"
  password = "ins3965!"
  url      = "https://10.23.248.120"
  insecure = true
}

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
}

resource "aci_l4_l7_service_graph_template" "service_graph" {
  tenant_dn = aci_tenant.tenant_for_contract.id
  name      = "test_service_graph"
}

resource "aci_l4_l7_service_graph_template" "service_graph2" {
  tenant_dn = aci_tenant.tenant_for_contract.id
  name      = "tf_service_graph"
}

// // apply_both_directions is selected [yes] by default and there is only one filter required
// resource "aci_contract_subject" "contract_subject" {
//   contract_dn   = aci_contract.democontract.id
//   name          = "contract_subject"
//   rev_flt_ports = "no"
// }aci_contract_subject.contract_subject.contract_dn

// apply_both_directions is not selected and there are two filters (consumer_to_provider and provider_to_consumer)
resource "aci_contract_subject" "contract_subject_2" {
  contract_dn   = aci_contract.democontract.id
  name          = "contract_subject_2"
  rev_flt_ports = "no"
  apply_both_directions = "yes"
  consumer_to_provider = {
    prio = "unspecified"
    target_dscp = "AF41"
    relation_vz_rs_in_term_graph_att = aci_l4_l7_service_graph_template.service_graph2.id
  }
  provider_to_consumer = {
    prio = "unspecified"
    target_dscp = "AF32"
    relation_vz_rs_out_term_graph_att = aci_l4_l7_service_graph_template.service_graph2.id
  }
}


// relation_vz_rs_out_term_graph_att = aci_l4_l7_service_graph_template.service_graph.id
// relation_vz_rs_filt_att = [{

//     },
//     {
      
//     }]
//   }

// data "aci_contract_subject" "example" {
//   contract_dn = aci_contract_subject.contract_subject_2.contract_dn
//   name          = aci_contract_subject.contract_subject_2.name
// }

// output "name" {
//   value = data.aci_contract_subject.example
// }


