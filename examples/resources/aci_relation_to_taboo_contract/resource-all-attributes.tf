
resource "aci_relation_to_taboo_contract" "full_example_application_epg" {
  parent_dn           = aci_application_epg.example.id
  annotation          = "annotation"
  taboo_contract_name = aci_taboo_contract.example.name
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
}
