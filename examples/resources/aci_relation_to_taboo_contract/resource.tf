
resource "aci_relation_to_taboo_contract" "example_application_epg" {
  parent_dn           = aci_application_epg.example.id
  taboo_contract_name = aci_taboo_contract.example.name
}
