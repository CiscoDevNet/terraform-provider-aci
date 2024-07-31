
data "aci_relation_to_imported_contract" "example_endpoint_security_group" {
  parent_dn              = aci_endpoint_security_group.example.id
  imported_contract_name = aci_imported_contract.example.name
}
