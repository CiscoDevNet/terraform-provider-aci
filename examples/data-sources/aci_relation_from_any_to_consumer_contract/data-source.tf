
data "aci_relation_from_any_to_consumer_contract" "example_any" {
  parent_dn     = aci_any.example.id
  contract_name = aci_contract.example.name
}
