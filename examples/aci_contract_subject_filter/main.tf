terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_tenant" "tenant_for_contract_filter" {
  name        = "tenant_for_contract_filter"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_contract" "test_contract" {
  tenant_dn                = aci_tenant.tenant_for_contract_filter.id
  name                     = "test_tf_contract"
}

resource "aci_contract_subject" "contract_subject" {
  contract_dn   = aci_contract.test_contract.id
  name          = "contract_subject"
  rev_flt_ports = "no"
}

resource "aci_contract_subject_filter" "contract_subject_filter" {
  contract_subject_dn = aci_contract_subject.contract_subject.id
  action = "permit"
  directives = ["log"]
  priority_override = "default"
  tn_vz_filter_name = "test_filter"
}

data "aci_contract_subject_filter" "example" {
  contract_subject_dn  = aci_contract_subject_filter.contract_subject_filter.contract_subject_dn
  tn_vz_filter_name  = aci_contract_subject_filter.contract_subject_filter.tn_vz_filter_name
}

output "name" {
  value = data.aci_contract_subject_filter.example
}
