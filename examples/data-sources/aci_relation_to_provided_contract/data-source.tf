
data "aci_relation_to_provided_contract" "example_application_epg" {
  parent_dn     = aci_application_epg.example.id
  contract_name = aci_contract.example.name
}

data "aci_relation_to_provided_contract" "example_endpoint_security_group" {
  parent_dn     = aci_endpoint_security_group.example.id
  contract_name = aci_contract.example.name
}
