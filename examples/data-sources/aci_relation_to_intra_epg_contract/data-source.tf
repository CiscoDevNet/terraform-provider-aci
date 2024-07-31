
data "aci_relation_to_intra_epg_contract" "example_endpoint_security_group" {
  parent_dn     = aci_endpoint_security_group.example.id
  contract_name = aci_contract.example.name
}
