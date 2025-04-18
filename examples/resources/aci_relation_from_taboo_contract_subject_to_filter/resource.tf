
resource "aci_relation_from_taboo_contract_subject_to_filter" "example_taboo_contract_subject" {
  parent_dn   = aci_taboo_contract_subject.example.id
  filter_name = aci_filter.example.name
}
