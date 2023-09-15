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

data "aci_tenant" "tenant_for_contract" {
  name = "tenant_for_contract"
}

data "aci_contract" "democontract" {
  tenant_dn = data.aci_tenant.tenant_for_contract.id
  name      = "test_tf_contract"
}

data "aci_contract_subject" "contract_subject" {
  contract_dn = data.aci_contract.democontract.id
  name        = "contract_subject"
}

data "aci_filter" "test_filter" {
  tenant_dn = data.aci_tenant.tenant_for_contract.id
  name      = "test_tf_filter"
}

# The aci_contract_subject_one_way_filter should only be used with datasources as it will create conflicts with the aci_contract_subject resource
resource "aci_contract_subject_one_way_filter" "contract_subject_filter" {
  contract_subject_dn = one(data.aci_contract_subject.contract_subject.consumer_to_provider).id
  filter_dn           = data.aci_filter.test_filter.id
}

resource "aci_contract_subject_one_way_filter" "contract_subject_filter2" {
  contract_subject_dn = one(data.aci_contract_subject.contract_subject.provider_to_consumer).id
  filter_dn           = data.aci_filter.test_filter.id
}