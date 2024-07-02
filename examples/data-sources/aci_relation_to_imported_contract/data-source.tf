
data "aci_relation_to_imported_contract" "example_application_epg" {
  parent_dn              = aci_application_epg.example.id
  imported_contract_name = aci_imported_contract.example.name
}

data "aci_relation_to_imported_contract" "example_endpoint_security_group" {
  parent_dn              = aci_endpoint_security_group.example.id
  imported_contract_name = aci_imported_contract.example.name
}
