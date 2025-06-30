
resource "aci_relation_from_any_to_contract_interface" "example_any" {
  parent_dn              = aci_any.example.id
  imported_contract_name = aci_imported_contract.example.name
}
